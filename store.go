package versionstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/goravel/framework/support/carbon"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == CONSTRUCTORS ============================================================

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	TableName          string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
	Logger             *slog.Logger
}

// NewStore creates a new session store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	store := &store{
		tableName:          opts.TableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.tableName == "" {
		return nil, errors.New("version store: tableName is required")
	}

	if store.db == nil {
		return nil, errors.New("version store: DB is required")
	}

	// Determine the database type
	store.dbDriverName = database.DatabaseType(store.db)

	// Set the default logger, if not provided
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}

	store.logger = opts.Logger

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}

// == TYPE ====================================================================

// Store defines a store
type store struct {
	tableName          string
	db                 *sql.DB
	dbDriverName       string
	logger             *slog.Logger
	automigrateEnabled bool
	debugEnabled       bool
}

var _ StoreInterface = (*store)(nil)

// AutoMigrate auto migrate
func (store *store) AutoMigrate() error {
	sqlStr := store.sqlTableCreate()

	_, err := store.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (store *store) EnableDebug(debug bool) {
	store.debugEnabled = debug
}

func (store *store) logSql(sqlStr string, sqlParams ...interface{}) {
	if store.debugEnabled {
		log.Println(sqlStr, sqlParams)
	}
}

func (store *store) VersionCount(options VersionQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.versionQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql(sqlStr, params...)

	mapped, err := database.SelectToMapString(context.TODO(), store.db, sqlStr, params...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *store) VersionCreate(version VersionInterface) error {
	if version == nil {
		return errors.New("version is nil")
	}

	if version.ID() == "" {
		return errors.New("version id should not be empty")
	}

	if version.EntityID() == "" {
		return errors.New("version entity id should not be empty")
	}

	if version.EntityType() == "" {
		return errors.New("version entity type should not be empty")
	}

	version.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	version.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := version.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.tableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql(sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	version.MarkAsNotDirty()

	return nil
}

func (store *store) VersionDelete(version VersionInterface) error {
	if version == nil {
		return errors.New("version is nil")
	}

	return store.VersionDeleteByID(version.ID())
}

func (store *store) VersionDeleteByID(id string) error {
	if id == "" {
		return errors.New("version id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.tableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql(sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *store) VersionFindByID(id string) (VersionInterface, error) {
	if id == "" {
		return nil, errors.New("version id is empty")
	}

	list, err := store.VersionList(NewVersionQuery().SetID(id).SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) VersionList(options VersionQueryInterface) ([]VersionInterface, error) {
	q, columns, err := store.versionQuery(options)

	if err != nil {
		return []VersionInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []VersionInterface{}, errSql
	}

	store.logSql(sqlStr, sqlParams...)

	modelMaps, err := database.SelectToMapString(context.Background(), store.db, sqlStr, sqlParams...)

	if err != nil {
		return []VersionInterface{}, err
	}

	list := []VersionInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewVersionFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) VersionSoftDelete(version VersionInterface) error {
	if version == nil {
		return errors.New("version is nil")
	}

	version.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.VersionUpdate(version)
}

func (store *store) VersionSoftDeleteByID(id string) error {
	version, err := store.VersionFindByID(id)

	if err != nil {
		return err
	}

	return store.VersionSoftDelete(version)
}

// VersionUpdate updates a version
//
// Note!! There is no reason to call this method other than marking
// the version as soft deleted
func (store *store) VersionUpdate(version VersionInterface) error {
	if version == nil {
		return errors.New("version is nil")
	}

	dataChanged := version.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.tableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(version.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql(sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	version.MarkAsNotDirty()

	return err
}

func (store *store) versionQuery(options VersionQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.tableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(cast.ToUint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(cast.ToUint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.HasSortOrder() {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted versions requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}
