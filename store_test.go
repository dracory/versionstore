package versionstore

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	_ = os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreVersionCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "version_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetEntityType("webpage").
		SetEntityID("1").
		SetContent("content1")

	ctx := database.Context(context.Background(), db)
	err = store.VersionCreate(ctx, version)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreVersionDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "version_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetEntityType("webpage").
		SetEntityID("1").
		SetContent("content1")

	ctx := database.Context(context.Background(), db)

	err = store.VersionCreate(ctx, version)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.VersionDelete(ctx, version)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	versionFound, errFind := store.VersionFindByID(ctx, version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if versionFound != nil {
		t.Fatal("Exam MUST be nil")
		return
	}
}

func TestStoreVersionDeleteByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "version_delete_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetEntityType("webpage").
		SetEntityID("1").
		SetContent("content1")
	ctx := database.Context(context.Background(), db)

	err = store.VersionCreate(ctx, version)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.VersionDeleteByID(ctx, version.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	versionFound, errFind := store.VersionFindByID(ctx, version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if versionFound != nil {
		t.Fatal("Exam MUST be nil")
		return
	}
}

func TestStoreVersionFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "shop_version_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetEntityType("discount").
		SetEntityID("1").
		SetContent("content1")

	ctx := database.Context(context.Background(), db)

	err = store.VersionCreate(ctx, version)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	versionFound, errFind := store.VersionFindByID(ctx, version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if versionFound == nil {
		t.Fatal("Version MUST NOT be nil")
		return
	}

	if versionFound.ID() != version.ID() {
		t.Fatal("Version ID MUST be equal. Expected: ", version.ID(), " Found: ", versionFound.ID())
	}

	if versionFound.EntityType() != version.EntityType() {
		t.Fatal("Version entity type MUST be equal. Expected: ", version.EntityType(), " Found: ", versionFound.EntityType())
	}

	if versionFound.EntityID() != version.EntityID() {
		t.Fatal("Version entity id MUST be equal. Expected: ", version.EntityID(), " Found: ", versionFound.EntityID())
	}

	if versionFound.Content() != version.Content() {
		t.Fatal("Version content MUST be equal. Expected: ", version.Content(), " Found: ", versionFound.Content())
	}

	if !strings.Contains(versionFound.CreatedAt(), version.CreatedAt()) {
		t.Fatal("Version created at MUST be equal. Expected: ", version.CreatedAt(), " Found: ", versionFound.CreatedAt())
	}

	if !strings.Contains(versionFound.SoftDeletedAt(), version.SoftDeletedAt()) {
		t.Fatal("Version soft deleted at MUST be equal. Expected: ", version.SoftDeletedAt(), " Found: ", versionFound.SoftDeletedAt())
	}
}

func TestStoreVersionSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "version_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetEntityType("webpage").
		SetEntityID("1").
		SetContent("content1")

	ctx := database.Context(context.Background(), db)

	err = store.VersionCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.VersionSoftDelete(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	versionFound, errFind := store.VersionFindByID(ctx, version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if versionFound != nil {
		t.Fatal("Version MUST be nil")
		return
	}

	versionList, errList := store.VersionList(ctx, NewVersionQuery().
		SetID(version.ID()).
		SetSoftDeletedIncluded(true))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
		return
	}

	if len(versionList) != 1 {
		t.Fatal("Version list MUST be 1")
		return
	}
}

func TestStoreVersionUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "shop_version_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetEntityType("discount").
		SetEntityID("1").
		SetContent("content1")

	ctx := database.Context(context.Background(), db)

	err = store.VersionCreate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	now := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	version.SetSoftDeletedAt(now)

	err = store.VersionUpdate(ctx, version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	versionList, errList := store.VersionList(ctx, NewVersionQuery().
		SetID(version.ID()).
		SetSoftDeletedIncluded(true))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
	}

	if len(versionList) < 1 {
		t.Fatal("Version list MUST NOT be 0")
	}

	if !strings.Contains(versionList[0].SoftDeletedAt(), now) {
		t.Fatal("Version soft deleted at MUST be equal. Expected: ", now, " Found: ", versionList[0].SoftDeletedAt())
	}
}
