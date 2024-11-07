package versionstore

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/goravel/framework/support/carbon"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

// Store defines a store
type store struct {
	tableName          string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

var _ StoreInterface = (*store)(nil)

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	TableName          string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new session store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	store := &store{
		tableName:          opts.TableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.tableName == "" {
		return nil, errors.New("log store: logTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("log store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate auto migrate
func (store *store) AutoMigrate() error {
	sql := store.sqlTableCreate()

	_, err := store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (store *store) EnableDebug(debug bool) {
	store.debugEnabled = debug
}

func (store *store) VersionCount(options VersionQueryInterface) (int64, error) {
	options.SetCountOnly(true)
	q := store.versionQuery(options)

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	mapped, err := db.SelectToMapString(sqlStr, params...)
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

	if version.Revision() < 1 {
		return errors.New("version revision should be greater than 0")
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

	if store.debugEnabled {
		log.Println(sqlStr)
	}

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

	if store.debugEnabled {
		log.Println(sqlStr)
	}

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
	q := store.versionQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []VersionInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	version.MarkAsNotDirty()

	return err
}

func (store *store) versionQuery(options VersionQueryInterface) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.tableName)

	if options.ID() != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	// if len(options.IDIn) > 0 {
	// 	q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	// }

	// if options.CreatedAtGte != "" && options.CreatedAtLte != "" {
	// 	q = q.Where(goqu.C(COLUMN_CREATED_AT).Between(exp.NewRangeVal(options.CreatedAtGte, options.CreatedAtLte)))
	// } else if options.CreatedAtGte != "" {
	// 	q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte))
	// } else if options.CreatedAtLte != "" {
	// 	q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte))
	// }

	if !options.CountOnly() {
		if options.Limit() > 0 {
			q = q.Limit(uint(options.Limit()))
		}

		if options.Offset() > 0 {
			q = q.Offset(uint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	if options.WithSoftDeleted() {
		return q
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted)
}
