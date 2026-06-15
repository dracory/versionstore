package versionstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	"github.com/dromara/carbon/v2"
)

// == CONSTRUCTORS ============================================================

// NewStoreOptions define the options for creating a new version store
type NewStoreOptions struct {
	TableName          string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
	Logger             *slog.Logger
}

// NewStore creates a new version store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.DB == nil {
		return nil, errors.New("version store: DB is required")
	}

	if opts.TableName == "" {
		return nil, errors.New("version store: tableName is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	store := &storeImplementation{
		tableName:          opts.TableName,
		db:                 neatDB,
		automigrateEnabled: opts.AutomigrateEnabled,
		debugEnabled:       opts.DebugEnabled,
		logger:             logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// == TYPE ====================================================================

// Store defines a store
type storeImplementation struct {
	tableName          string
	db                 *neat.Database
	logger             *slog.Logger
	automigrateEnabled bool
	debugEnabled       bool
}

var _ StoreInterface = (*storeImplementation)(nil)

// AutoMigrate auto migrate (deprecated - use MigrateUp)
func (store *storeImplementation) AutoMigrate() error {
	return store.MigrateUp(context.Background())
}

// MigrateUp creates the table
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if store.db.Schema().HasTable(store.tableName) {
		if store.debugEnabled {
			store.logger.Info("MigrateUp: table already exists", "table", store.tableName)
		}
		return nil
	}

	err := store.db.Schema().Create(store.tableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 21)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_ENTITY_TYPE, 40)
		table.String(COLUMN_ENTITY_ID, 40)
		table.Text(COLUMN_CONTENT)
		table.DateTime(COLUMN_CREATED_AT)
		table.DateTime(COLUMN_SOFT_DELETED_AT)
	})

	if err != nil {
		if store.debugEnabled {
			store.logger.Error("MigrateUp failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the table
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if !store.db.Schema().HasTable(store.tableName) {
		if store.debugEnabled {
			store.logger.Info("MigrateDown: table does not exist", "table", store.tableName)
		}
		return nil
	}

	err := store.db.Schema().Drop(store.tableName)
	if err != nil {
		if store.debugEnabled {
			store.logger.Error("MigrateDown failed", "error", err)
		}
		return err
	}
	return nil
}

// EnableDebug - enables the debug option
func (store *storeImplementation) EnableDebug(debug bool) {
	store.debugEnabled = debug
	if debug {
		store.db.EnableDebug()
	} else {
		store.db.DisableDebug()
	}
}

// GetTableName returns the table name
func (store *storeImplementation) GetTableName() string {
	return store.tableName
}

// SetTableName sets the table name
func (store *storeImplementation) SetTableName(tableName string) {
	store.tableName = tableName
}

// VersionCount returns the count of versions matching the query options
func (store *storeImplementation) VersionCount(ctx context.Context, options VersionQueryInterface) (int64, error) {
	if ctx == nil {
		return 0, errors.New("ctx is nil")
	}

	q := store.buildQuery(options)

	var count int64
	err := q.Table(store.tableName).Count(&count)
	return count, err
}

// VersionCreate creates a new version
func (store *storeImplementation) VersionCreate(ctx context.Context, version VersionInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version store: version cannot be nil")
	}
	if version.ID() == "" {
		return errors.New("version store: version id should not be empty")
	}
	if version.EntityType() == "" {
		return errors.New("version store: version entity type should not be empty")
	}
	if version.EntityID() == "" {
		return errors.New("version store: version entity id should not be empty")
	}
	if version.GetCreatedAt() == "" {
		version.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}
	if version.GetSoftDeletedAt() == "" {
		version.SetSoftDeletedAt(MAX_DATETIME)
	}

	row := map[string]any{
		COLUMN_ID:              version.ID(),
		COLUMN_ENTITY_TYPE:     version.EntityType(),
		COLUMN_ENTITY_ID:       version.EntityID(),
		COLUMN_CONTENT:         version.Content(),
		COLUMN_CREATED_AT:      version.GetCreatedAtCarbon().StdTime(),
		COLUMN_SOFT_DELETED_AT: version.GetSoftDeletedAtCarbon().StdTime(),
	}

	return store.db.Query().Table(store.tableName).Create(row)
}

// VersionDelete deletes a version permanently
func (store *storeImplementation) VersionDelete(ctx context.Context, version VersionInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version is nil")
	}

	return store.VersionDeleteByID(ctx, version.ID())
}

// VersionDeleteByID deletes a version by ID permanently
func (store *storeImplementation) VersionDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("version id is empty")
	}

	_, err := store.db.Query().
		Table(store.tableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()
	return err
}

// VersionFindByID finds a version by ID
func (store *storeImplementation) VersionFindByID(ctx context.Context, id string) (VersionInterface, error) {
	if id == "" {
		return nil, errors.New("version store: version id is required")
	}

	list, err := store.VersionList(ctx, NewVersionQuery().SetID(id).SetLimit(1))
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// VersionList returns a list of versions matching the query options
func (store *storeImplementation) VersionList(ctx context.Context, options VersionQueryInterface) ([]VersionInterface, error) {
	if ctx == nil {
		return nil, errors.New("ctx is nil")
	}

	type versionRow struct {
		ID            string    `db:"id"`
		EntityType    string    `db:"entity_type"`
		EntityID      string    `db:"entity_id"`
		Content       string    `db:"content"`
		CreatedAt     time.Time `db:"created_at"`
		SoftDeletedAt time.Time `db:"soft_deleted_at"`
	}

	q := store.buildQuery(options)
	q = q.Table(store.tableName)

	if len(options.Columns()) > 0 {
		q = q.Select(options.Columns())
	}

	var rows []versionRow
	if err := q.Get(&rows); err != nil {
		return []VersionInterface{}, err
	}

	list := make([]VersionInterface, 0, len(rows))
	for _, r := range rows {
		v := &version{}
		v.SetID(r.ID)
		v.SetEntityType(r.EntityType)
		v.SetEntityID(r.EntityID)
		v.SetContent(r.Content)
		v.CreatedAtField.CreatedAt = r.CreatedAt
		v.SoftDeletedAt = r.SoftDeletedAt
		list = append(list, v)
	}

	return list, nil
}

// VersionSoftDelete soft deletes a version
func (store *storeImplementation) VersionSoftDelete(ctx context.Context, version VersionInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version is nil")
	}

	version.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.VersionUpdate(ctx, version)
}

// VersionSoftDeleteByID soft deletes a version by ID
func (store *storeImplementation) VersionSoftDeleteByID(ctx context.Context, id string) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if id == "" {
		return errors.New("version id is empty")
	}

	version, err := store.VersionFindByID(ctx, id)
	if err != nil {
		return err
	}
	if version == nil {
		return errors.New("version not found")
	}

	return store.VersionSoftDelete(ctx, version)
}

// VersionUpdate updates a version
//
// Note!! There is no reason to call this method other than marking
// the version as soft deleted
func (store *storeImplementation) VersionUpdate(ctx context.Context, version VersionInterface) error {
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	if version == nil {
		return errors.New("version is nil")
	}

	row := map[string]any{
		COLUMN_SOFT_DELETED_AT: version.GetSoftDeletedAtCarbon().StdTime(),
	}

	_, err := store.db.Query().Table(store.tableName).Where(COLUMN_ID+" = ?", version.ID()).Update(row)
	return err
}

// == QUERY BUILDER ==========================================================

// buildQuery builds a neat query from the version query interface.
func (store *storeImplementation) buildQuery(options VersionQueryInterface) contractsorm.Query {
	q := store.db.Query()

	if options == nil {
		return q
	}

	if options.HasID() && options.ID() != "" {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasEntityType() && options.EntityType() != "" {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType())
	}

	if options.HasEntityID() && options.EntityID() != "" {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
	}

	if options.HasLimit() && options.Limit() > 0 {
		q = q.Limit(options.Limit())
	}

	if options.HasOffset() && options.Offset() > 0 {
		q = q.Offset(int(options.Offset()))
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if options.HasSortOrder() && options.SortOrder() == "asc" {
			q = q.OrderBy(options.OrderBy())
		} else {
			q = q.OrderByDesc(options.OrderBy())
		}
	}

	if options.HasSoftDeletedIncluded() && options.SoftDeletedIncluded() {
		q = q.WithSoftDeleted()
	} else {
		q = q.Where(COLUMN_SOFT_DELETED_AT+" = ?", carbon.Parse(MAX_DATETIME, carbon.UTC).StdTime())
	}

	return q
}
