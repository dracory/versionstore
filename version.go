package versionstore

import (
	"time"

	"github.com/dracory/neat/database/orm"
	"github.com/dracory/neat/database/soft_delete"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == CONSTRUCTOR =============================================================

// NewVersion creates a new version with a generated ID and current timestamp
func NewVersion() VersionInterface {
	o := &version{}
	o.SetID(neatuid.GenerateShortID())
	o.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	o.SetSoftDeletedAt(MAX_DATETIME)
	return o
}

// NewVersionFromExistingData creates a version from existing data
func NewVersionFromExistingData(data map[string]string) VersionInterface {
	o := &version{}
	o.SetID(data[COLUMN_ID])
	o.SetEntityType(data[COLUMN_ENTITY_TYPE])
	o.SetEntityID(data[COLUMN_ENTITY_ID])
	o.SetContent(data[COLUMN_CONTENT])
	if v, ok := data[COLUMN_CREATED_AT]; ok {
		o.SetCreatedAt(v)
	}
	if v, ok := data[COLUMN_SOFT_DELETED_AT]; ok {
		o.SetSoftDeletedAt(v)
	}
	return o
}

// == CLASS ==================================================================

var _ VersionInterface = (*version)(nil)

type version struct {
	orm.ShortID

	EntityTypeField string `db:"entity_type"`
	EntityIDField   string `db:"entity_id"`
	ContentField    string `db:"content"`

	CreatedAtField orm.CreatedAt
	soft_delete.SoftDeletesMaxDate
}

// == METHODS =================================================================

// IsSoftDeleted returns true if the version is soft deleted.
func (o *version) IsSoftDeleted() bool {
	return o.SoftDeletedAt.Before(time.Now().UTC())
}

// == SETTERS AND GETTERS =====================================================

// ID returns the id of the version.
func (o *version) ID() string {
	return o.ShortID.ID
}

// SetID sets the id of the version.
func (o *version) SetID(id string) VersionInterface {
	o.ShortID.ID = id
	return o
}

// EntityType returns the entity type of the version.
func (o *version) EntityType() string {
	return o.EntityTypeField
}

// SetEntityType sets the entity type of the version.
func (o *version) SetEntityType(entityType string) VersionInterface {
	o.EntityTypeField = entityType
	return o
}

// EntityID returns the entity id of the version.
func (o *version) EntityID() string {
	return o.EntityIDField
}

// SetEntityID sets the entity id of the version.
func (o *version) SetEntityID(entityID string) VersionInterface {
	o.EntityIDField = entityID
	return o
}

// Content returns the content of the version.
func (o *version) Content() string {
	return o.ContentField
}

// SetContent sets the content of the version.
func (o *version) SetContent(content string) VersionInterface {
	o.ContentField = content
	return o
}

// GetCreatedAt returns the created at time of the version.
func (o *version) GetCreatedAt() string {
	if o.CreatedAtField.CreatedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt).ToDateTimeString()
}

// GetCreatedAtCarbon returns the created at time of the version as a carbon object.
func (o *version) GetCreatedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.CreatedAtField.CreatedAt)
}

// SetCreatedAt sets the created at time of the version.
func (o *version) SetCreatedAt(createdAt string) VersionInterface {
	if createdAt == "" {
		return o
	}
	o.CreatedAtField.CreatedAt = carbon.Parse(createdAt, carbon.UTC).StdTime()
	return o
}

// GetSoftDeletedAt returns the soft deleted at time of the version.
func (o *version) GetSoftDeletedAt() string {
	if o.SoftDeletedAt.IsZero() {
		return ""
	}
	return carbon.CreateFromStdTime(o.SoftDeletedAt).ToDateTimeString()
}

// GetSoftDeletedAtCarbon returns the soft deleted at time of the version as a carbon object.
func (o *version) GetSoftDeletedAtCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(o.SoftDeletedAt)
}

// SetSoftDeletedAt sets the soft deleted at time of the version.
func (o *version) SetSoftDeletedAt(softDeletedAt string) VersionInterface {
	if softDeletedAt == "" {
		return o
	}
	o.SoftDeletedAt = carbon.Parse(softDeletedAt, carbon.UTC).StdTime()
	return o
}
