package versionstore

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
	Revision() int
	SetRevision(revision int) VersionInterface
	Content() string
	SetContent(content string) VersionInterface
	CreatedAt() string
	SetCreatedAt(createdAt string) VersionInterface
	SoftDeletedAt() string
	SetSoftDeletedAt(softDeletedAt string) VersionInterface
}
