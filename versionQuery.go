package versionstore

import "errors"

// NewVersionQuery creates a new version query
func NewVersionQuery() VersionQueryInterface {
	return &versionQuery{
		properties: map[string]any{},
	}
}

// versionQuery implements VersionQueryInterface
type versionQuery struct {
	properties map[string]any
}

// Validate validates the query parameters
func (q *versionQuery) Validate() error {
	if q.HasEntityID() && q.EntityID() == "" {
		return errors.New("version query. entity_id cannot be empty")
	}

	if q.HasEntityType() && q.EntityType() == "" {
		return errors.New("version query. entity_type cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("version query. id cannot be empty")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("version query. limit cannot be negative")
	}

	if q.HasLimit() && q.Limit() < 1 {
		return errors.New("version query. limit cannot be less than 1")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("version query. offset cannot be negative")
	}

	return nil
}

// Columns returns the columns to select
func (q *versionQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

// SetColumns sets the columns to select
func (q *versionQuery) SetColumns(columns []string) VersionQueryInterface {
	q.properties["columns"] = columns
	return q
}

// IsCountOnly returns true if only count is requested
func (q *versionQuery) IsCountOnly() bool {
	if !q.hasProperty("count_only") {
		return false
	}

	return q.properties["count_only"].(bool)
}

// HasCountOnly returns true if count_only is set
func (q *versionQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

// SetCountOnly sets whether to return only count
func (q *versionQuery) SetCountOnly(countOnly bool) VersionQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

// HasEntityID returns true if entity_id is set
func (q *versionQuery) HasEntityID() bool {
	return q.hasProperty("entity_id")
}

// EntityID returns the entity ID
func (q *versionQuery) EntityID() string {
	if !q.hasProperty("entity_id") {
		return ""
	}

	return q.properties["entity_id"].(string)
}

// SetEntityID sets the entity ID
func (q *versionQuery) SetEntityID(entityID string) VersionQueryInterface {
	q.properties["entity_id"] = entityID
	return q
}

// HasEntityType returns true if entity_type is set
func (q *versionQuery) HasEntityType() bool {
	return q.hasProperty("entity_type")
}

// EntityType returns the entity type
func (q *versionQuery) EntityType() string {
	if !q.hasProperty("entity_type") {
		return ""
	}

	return q.properties["entity_type"].(string)
}

// SetEntityType sets the entity type
func (q *versionQuery) SetEntityType(entityType string) VersionQueryInterface {
	q.properties["entity_type"] = entityType
	return q
}

// HasID returns true if id is set
func (q *versionQuery) HasID() bool {
	return q.hasProperty("id")
}

// ID returns the version ID
func (q *versionQuery) ID() string {
	if !q.hasProperty("id") {
		return ""
	}

	return q.properties["id"].(string)
}

// SetID sets the version ID
func (q *versionQuery) SetID(id string) VersionQueryInterface {
	q.properties["id"] = id
	return q
}

// HasLimit returns true if limit is set
func (q *versionQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

// Limit returns the query limit
func (q *versionQuery) Limit() int {
	if !q.hasProperty("limit") {
		return 0
	}

	return q.properties["limit"].(int)
}

// SetLimit sets the query limit
func (q *versionQuery) SetLimit(limit int) VersionQueryInterface {
	q.properties["limit"] = limit
	return q
}

// HasOffset returns true if offset is set
func (q *versionQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

// Offset returns the query offset
func (q *versionQuery) Offset() int64 {
	if !q.hasProperty("offset") {
		return 0
	}

	return q.properties["offset"].(int64)
}

// SetOffset sets the query offset
func (q *versionQuery) SetOffset(offset int64) VersionQueryInterface {
	q.properties["offset"] = offset
	return q
}

// HasOrderBy returns true if order_by is set
func (q *versionQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

// OrderBy returns the order by field
func (q *versionQuery) OrderBy() string {
	if !q.hasProperty("order_by") {
		return ""
	}

	return q.properties["order_by"].(string)
}

// SetOrderBy sets the order by field
func (q *versionQuery) SetOrderBy(orderBy string) VersionQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

// HasSortOrder returns true if sort_order is set
func (q *versionQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

// SortOrder returns the sort order (ASC or DESC)
func (q *versionQuery) SortOrder() string {
	if !q.hasProperty("sort_order") {
		return ""
	}

	return q.properties["sort_order"].(string)
}

// SetSortOrder sets the sort order (ASC or DESC)
func (q *versionQuery) SetSortOrder(sortOrder string) VersionQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

// HasSoftDeletedIncluded returns true if soft_deleted_included is set
func (q *versionQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

// SoftDeletedIncluded returns true if soft deleted versions should be included
func (q *versionQuery) SoftDeletedIncluded() bool {
	if q.hasProperty("soft_deleted_included") {
		return q.properties["soft_deleted_included"].(bool)
	}

	return false
}

// SetSoftDeletedIncluded sets whether to include soft deleted versions
func (q *versionQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) VersionQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

// hasProperty returns true if the property exists in the map
func (q *versionQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
