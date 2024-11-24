package versionstore

import "errors"

func NewVersionQuery() VersionQueryInterface {
	return &versionQuery{
		properties: map[string]any{},
	}
}

type versionQuery struct {
	properties map[string]any
}

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

func (q *versionQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *versionQuery) SetColumns(columns []string) VersionQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *versionQuery) IsCountOnly() bool {
	if !q.hasProperty("count_only") {
		return false
	}

	return q.properties["count_only"].(bool)
}

func (q *versionQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *versionQuery) SetCountOnly(countOnly bool) VersionQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *versionQuery) HasEntityID() bool {
	return q.hasProperty("entity_id")
}

func (q *versionQuery) EntityID() string {
	if !q.hasProperty("entity_id") {
		return ""
	}

	return q.properties["entity_id"].(string)
}

func (q *versionQuery) SetEntityID(entityID string) VersionQueryInterface {
	q.properties["entity_id"] = entityID
	return q
}

func (q *versionQuery) HasEntityType() bool {
	return q.hasProperty("entity_type")
}

func (q *versionQuery) EntityType() string {
	if !q.hasProperty("entity_type") {
		return ""
	}

	return q.properties["entity_type"].(string)
}

func (q *versionQuery) SetEntityType(entityType string) VersionQueryInterface {
	q.properties["entity_type"] = entityType
	return q
}

func (q *versionQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *versionQuery) ID() string {
	if !q.hasProperty("id") {
		return ""
	}

	return q.properties["id"].(string)
}

func (q *versionQuery) SetID(id string) VersionQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *versionQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *versionQuery) Limit() int {
	if !q.hasProperty("limit") {
		return 0
	}

	return q.properties["limit"].(int)
}

func (q *versionQuery) SetLimit(limit int) VersionQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *versionQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *versionQuery) Offset() int64 {
	if !q.hasProperty("offset") {
		return 0
	}

	return q.properties["offset"].(int64)
}

func (q *versionQuery) SetOffset(offset int64) VersionQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *versionQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *versionQuery) OrderBy() string {
	if !q.hasProperty("order_by") {
		return ""
	}

	return q.properties["order_by"].(string)
}

func (q *versionQuery) SetOrderBy(orderBy string) VersionQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *versionQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *versionQuery) SortOrder() string {
	if !q.hasProperty("sort_order") {
		return ""
	}

	return q.properties["sort_order"].(string)
}

func (q *versionQuery) SetSortOrder(sortOrder string) VersionQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *versionQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

func (q *versionQuery) SoftDeletedIncluded() bool {
	if q.hasProperty("soft_deleted_included") {
		return q.properties["soft_deleted_included"].(bool)
	}

	return false
}

func (q *versionQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) VersionQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *versionQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
