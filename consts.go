package versionstore

// Column names for the version table
const (
	COLUMN_CONTENT         = "content"
	COLUMN_CREATED_AT      = "created_at"
	COLUMN_ENTITY_ID       = "entity_id"
	COLUMN_ENTITY_TYPE     = "entity_type"
	COLUMN_ID              = "id"
	COLUMN_SOFT_DELETED_AT = "soft_deleted_at"
)

// MAX_DATETIME is a far-future datetime used as the default soft-delete sentinel.
const MAX_DATETIME = "9999-12-31 23:59:59"
