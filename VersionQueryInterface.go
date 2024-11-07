package versionstore

type VersionQueryInterface interface {
	ID() string
	SetID(id string) VersionQueryInterface
	EntityType() string
	SetEntityID(entityID string) VersionQueryInterface
	EntityID() string
	SetEntityType(entityType string) VersionQueryInterface
	Offset() int64
	SetOffset(offset int64) VersionQueryInterface
	Limit() int
	SetLimit(limit int) VersionQueryInterface
	SortOrder() string
	SetSortOrder(sortOrder string) VersionQueryInterface
	OrderBy() string
	SetOrderBy(orderBy string) VersionQueryInterface
	CountOnly() bool
	SetCountOnly(countOnly bool) VersionQueryInterface
	WithSoftDeleted() bool
	SetWithSoftDeleted(withSoftDeleted bool) VersionQueryInterface
}
