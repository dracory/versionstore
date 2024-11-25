package versionstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

// == CONSTRUCTOR =============================================================

func NewVersion() VersionInterface {
	return &version{
		id:            uid.NanoUid(),
		createdAt:     carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC),
		softDeletedAt: sb.MAX_DATETIME,
	}
}

func NewVersionFromExistingData(data map[string]string) VersionInterface {
	return &version{
		id:            data[COLUMN_ID],
		entityType:    data[COLUMN_ENTITY_TYPE],
		entityID:      data[COLUMN_ENTITY_ID],
		content:       data[COLUMN_CONTENT],
		createdAt:     data[COLUMN_CREATED_AT],
		softDeletedAt: data[COLUMN_SOFT_DELETED_AT],
	}
}

// == CLASS ==================================================================

type version struct {
	id            string
	entityType    string
	entityID      string
	content       string
	createdAt     string
	softDeletedAt string
}

var _ VersionInterface = (*version)(nil)

func (v *version) Data() map[string]string {
	return map[string]string{
		COLUMN_ID:              v.id,
		COLUMN_ENTITY_TYPE:     v.entityType,
		COLUMN_ENTITY_ID:       v.entityID,
		COLUMN_CONTENT:         v.content,
		COLUMN_CREATED_AT:      v.createdAt,
		COLUMN_SOFT_DELETED_AT: v.softDeletedAt,
	}
}

// DataChanged returns the data that has changed
//
// This entity does not have many fields that can be changed
// only the soft deleted field can be changed, the other fields
// should not be changed really
func (v *version) DataChanged() map[string]string {
	return map[string]string{
		COLUMN_SOFT_DELETED_AT: v.softDeletedAt,
	}
}

func (v *version) MarkAsNotDirty() {
	// there is nothing to do here
}

func (v *version) Content() string {
	return v.content
}

func (v *version) SetContent(content string) VersionInterface {
	v.content = content
	return v
}

func (v *version) CreatedAt() string {
	return v.createdAt
}

func (v *version) SetCreatedAt(createdAt string) VersionInterface {
	v.createdAt = createdAt
	return v
}

func (v *version) EntityType() string {
	return v.entityType
}

func (v *version) SetEntityType(entityType string) VersionInterface {
	v.entityType = entityType
	return v
}

func (v *version) EntityID() string {
	return v.entityID
}

func (v *version) SetEntityID(entityID string) VersionInterface {
	v.entityID = entityID
	return v
}

func (v *version) ID() string {
	return v.id
}

func (v *version) SetID(id string) VersionInterface {
	v.id = id
	return v
}

func (v *version) SoftDeletedAt() string {
	return v.softDeletedAt
}

func (v *version) SetSoftDeletedAt(softDeletedAt string) VersionInterface {
	v.softDeletedAt = softDeletedAt
	return v
}
