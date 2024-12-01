package versionstore

import "context"

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
	VersionCreate(context context.Context, version VersionInterface) error
	VersionFindByID(context context.Context, versionID string) (VersionInterface, error)
	VersionList(context context.Context, query VersionQueryInterface) ([]VersionInterface, error)
	VersionUpdate(context context.Context, version VersionInterface) error
	VersionDelete(context context.Context, version VersionInterface) error
	VersionDeleteByID(context context.Context, versionID string) error
	VersionSoftDelete(context context.Context, version VersionInterface) error
	VersionSoftDeleteByID(context context.Context, versionID string) error
}

type VersionInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	ID() string
	SetID(id string) VersionInterface
	EntityType() string
	SetEntityType(entityType string) VersionInterface
	EntityID() string
	SetEntityID(entityID string) VersionInterface
	Content() string
	SetContent(content string) VersionInterface
	CreatedAt() string
	SetCreatedAt(createdAt string) VersionInterface
	SoftDeletedAt() string
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
