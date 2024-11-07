package versionstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
	VersionCreate(version VersionInterface) error
	VersionFindByID(versionID string) (VersionInterface, error)
	VersionList(query VersionQueryInterface) ([]VersionInterface, error)
	VersionUpdate(version VersionInterface) error
	VersionDelete(version VersionInterface) error
	VersionDeleteByID(versionID string) error
	VersionSoftDelete(version VersionInterface) error
	VersionSoftDeleteByID(versionID string) error
}
