package versionstore

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dracory/database"
	"github.com/dromara/carbon/v2"
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

func TestStoreVersionSoftDeleteByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "version_soft_delete_by_id",
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

	err = store.VersionSoftDeleteByID(ctx, version.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	// After soft delete, VersionFindByID should NOT find the record (soft deleted excluded by default)
	versionFound, errFind := store.VersionFindByID(ctx, version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if versionFound != nil {
		t.Fatal("Version MUST be nil after soft delete (FindByID excludes soft deleted)")
		return
	}

	// But the record should still exist if we include soft deleted
	versionList, errList := store.VersionList(ctx, NewVersionQuery().
		SetID(version.ID()).
		SetSoftDeletedIncluded(true))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
		return
	}

	if len(versionList) != 1 {
		t.Fatal("Version list MUST be 1, got:", len(versionList))
		return
	}

	if versionList[0].SoftDeletedAt() == "" {
		t.Fatal("Version SoftDeletedAt MUST NOT be empty after soft delete by ID")
	}
}

func TestStoreVersionList_Ordering(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "version_list_ordering",
		AutomigrateEnabled: true,
		DebugEnabled:       true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	ctx := database.Context(context.Background(), db)
	entityID := "entity-123"
	entityType := "webpage"

	// Create first version
	version1 := NewVersion().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetContent("content1")

	err = store.VersionCreate(ctx, version1)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	// Delay to ensure different timestamps (SQLite datetime has second precision)
	time.Sleep(1 * time.Second)

	// Create second version
	version2 := NewVersion().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetContent("content2")

	err = store.VersionCreate(ctx, version2)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	// List with DESC order - should return newest first
	versionList, errList := store.VersionList(ctx, NewVersionQuery().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetOrderBy(COLUMN_CREATED_AT).
		SetSortOrder("desc"))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
		return
	}

	if len(versionList) != 2 {
		t.Fatal("Version list MUST be 2, got:", len(versionList))
		return
	}

	// DESC order should return content2 first (newest)
	if versionList[0].Content() != "content2" {
		t.Fatal("First version MUST be 'content2' (newest) with DESC order. Got:", versionList[0].Content())
	}

	if versionList[1].Content() != "content1" {
		t.Fatal("Second version MUST be 'content1' (oldest) with DESC order. Got:", versionList[1].Content())
	}
}
