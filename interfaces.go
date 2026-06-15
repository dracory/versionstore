package versionstore

import (
	"context"
	"database/sql"

	"github.com/dromara/carbon/v2"
)

type StoreInterface interface {
	// GetTableName returns the table name
	GetTableName() string
	// SetTableName sets the table name
	SetTableName(tableName string)

	// MigrateDown drops the table
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error
	// MigrateUp creates the table
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	EnableDebug(debug bool)
	VersionCreate(ctx context.Context, version VersionInterface) error
	VersionFindByID(ctx context.Context, versionID string) (VersionInterface, error)
	VersionList(ctx context.Context, query VersionQueryInterface) ([]VersionInterface, error)
	VersionUpdate(ctx context.Context, version VersionInterface) error
	VersionDelete(ctx context.Context, version VersionInterface) error
	VersionDeleteByID(ctx context.Context, versionID string) error
	VersionSoftDelete(ctx context.Context, version VersionInterface) error
	VersionSoftDeleteByID(ctx context.Context, versionID string) error
}

type VersionInterface interface {
	IsSoftDeleted() bool

	ID() string
	SetID(id string) VersionInterface

	EntityType() string
	SetEntityType(entityType string) VersionInterface

	EntityID() string
	SetEntityID(entityID string) VersionInterface

	Content() string
	SetContent(content string) VersionInterface

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) VersionInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(softDeletedAt string) VersionInterface
}

type VersionQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) VersionQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) VersionQueryInterface

	HasID() bool
	ID() string
	SetID(id string) VersionQueryInterface

	HasEntityID() bool
	EntityID() string
	SetEntityID(entityID string) VersionQueryInterface

	HasEntityType() bool
	EntityType() string
	SetEntityType(entityType string) VersionQueryInterface

	HasOffset() bool
	Offset() int64
	SetOffset(offset int64) VersionQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) VersionQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) VersionQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) VersionQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(includeSoftDeleted bool) VersionQueryInterface
}
