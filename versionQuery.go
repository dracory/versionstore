package versionstore

type versionQuery struct {
	id              string
	entityType      string
	entityID        string
	countOnly       bool
	offset          int64
	limit           int
	sortOrder       string
	orderBy         string
	withSoftDeleted bool
}

func NewVersionQuery() *versionQuery {
	return &versionQuery{}
}

var _ VersionQueryInterface = &versionQuery{}

func (q *versionQuery) CountOnly() bool {
	return q.countOnly
}

func (q *versionQuery) SetCountOnly(countOnly bool) VersionQueryInterface {
	q.countOnly = countOnly
	return q
}

func (q *versionQuery) EntityID() string {
	return q.entityID
}

func (q *versionQuery) SetEntityID(entityID string) VersionQueryInterface {
	q.entityID = entityID
	return q
}

func (q *versionQuery) EntityType() string {
	return q.entityType
}

func (q *versionQuery) SetEntityType(entityType string) VersionQueryInterface {
	q.entityType = entityType
	return q
}

func (q *versionQuery) ID() string {
	return q.id
}

func (q *versionQuery) SetID(id string) VersionQueryInterface {
	q.id = id
	return q
}

func (q *versionQuery) Limit() int {
	return q.limit
}

func (q *versionQuery) SetLimit(limit int) VersionQueryInterface {
	q.limit = limit
	return q
}

func (q *versionQuery) Offset() int64 {
	return q.offset
}

func (q *versionQuery) SetOffset(offset int64) VersionQueryInterface {
	q.offset = offset
	return q
}

func (q *versionQuery) OrderBy() string {
	return q.orderBy
}

func (q *versionQuery) SetOrderBy(orderBy string) VersionQueryInterface {
	q.orderBy = orderBy
	return q
}

func (q *versionQuery) SortOrder() string {
	return q.sortOrder
}

func (q *versionQuery) SetSortOrder(sortOrder string) VersionQueryInterface {
	q.sortOrder = sortOrder
	return q
}

func (q *versionQuery) WithSoftDeleted() bool {
	return q.withSoftDeleted
}

func (q *versionQuery) SetWithSoftDeleted(withSoftDeleted bool) VersionQueryInterface {
	q.withSoftDeleted = withSoftDeleted
	return q
}
