package versionstore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
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
		SetRevision(1).
		SetContent("content1")

	err = store.VersionCreate(version)

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
		SetRevision(1).
		SetContent("content1")

	err = store.VersionCreate(version)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.VersionDelete(version)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	versionFound, errFind := store.VersionFindByID(version.ID())

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
		SetRevision(1).
		SetContent("content1")

	err = store.VersionCreate(version)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.VersionDeleteByID(version.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	versionFound, errFind := store.VersionFindByID(version.ID())

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
		SetRevision(1).
		SetContent("content1")

	err = store.VersionCreate(version)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	versionFound, errFind := store.VersionFindByID(version.ID())

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

	if versionFound.Revision() != version.Revision() {
		t.Fatal("Version revision MUST be equal. Expected: ", version.Revision(), " Found: ", versionFound.Revision())
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
		SetRevision(1).
		SetContent("content1")

	err = store.VersionCreate(version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	err = store.VersionSoftDelete(version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	versionFound, errFind := store.VersionFindByID(version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if versionFound != nil {
		t.Fatal("Version MUST be nil")
		return
	}

	versionList, errList := store.VersionList(NewVersionQuery().
		SetID(version.ID()).
		SetWithSoftDeleted(true))

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
		DB:                     db,
		TableName:              "shop_version_update",
		OrderTableName:         "shop_order_update",
		OrderLineItemTableName: "shop_order_line_item_update",
		ProductTableName:       "shop_product_update",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	version := NewVersion().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE").
		SetDescription("DISCOUNT_DESCRIPTION").
		SetType(DISCOUNT_TYPE_AMOUNT).
		SetAmount(19.99).
		SetStartsAt(`2022-01-01 00:00:00`).
		SetEndsAt(`2022-01-01 23:59:59`)

	err = store.VersionCreate(version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	version.SetTitle("DISCOUNT_TITLE_UPDATED")

	err = store.VersionUpdate(version)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	versionFound, errFind := store.VersionFindByID(version.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if versionFound == nil {
		t.Fatal("Version MUST NOT be nil")
	}

	if versionFound.Title() != "DISCOUNT_TITLE_UPDATED" {
		t.Fatal("Version title MUST BE 'DISCOUNT_TITLE_UPDATED', found: ", versionFound.Title())
	}
}
