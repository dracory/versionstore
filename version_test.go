package versionstore

import (
	"testing"

	"github.com/dracory/sb"
)

func TestNewVersion(t *testing.T) {
	version := NewVersion()

	if version == nil {
		t.Fatal("NewVersion() returned nil")
	}

	if version.ID() == "" {
		t.Error("NewVersion() should generate an ID")
	}

	if version.CreatedAt() == "" {
		t.Error("NewVersion() should set CreatedAt")
	}

	if version.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("NewVersion() SoftDeletedAt should be %s, got %s", sb.MAX_DATETIME, version.SoftDeletedAt())
	}
}

func TestNewVersionFromExistingData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "test-id",
		COLUMN_ENTITY_TYPE:     "test-entity-type",
		COLUMN_ENTITY_ID:       "test-entity-id",
		COLUMN_CONTENT:         "test-content",
		COLUMN_CREATED_AT:      "2024-01-01 00:00:00",
		COLUMN_SOFT_DELETED_AT: sb.MAX_DATETIME,
	}

	version := NewVersionFromExistingData(data)

	if version == nil {
		t.Fatal("NewVersionFromExistingData() returned nil")
	}

	if version.ID() != "test-id" {
		t.Errorf("ID() = %s, want %s", version.ID(), "test-id")
	}

	if version.EntityType() != "test-entity-type" {
		t.Errorf("EntityType() = %s, want %s", version.EntityType(), "test-entity-type")
	}

	if version.EntityID() != "test-entity-id" {
		t.Errorf("EntityID() = %s, want %s", version.EntityID(), "test-entity-id")
	}

	if version.Content() != "test-content" {
		t.Errorf("Content() = %s, want %s", version.Content(), "test-content")
	}

	if version.CreatedAt() != "2024-01-01 00:00:00" {
		t.Errorf("CreatedAt() = %s, want %s", version.CreatedAt(), "2024-01-01 00:00:00")
	}

	if version.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("SoftDeletedAt() = %s, want %s", version.SoftDeletedAt(), sb.MAX_DATETIME)
	}
}

func TestVersionData(t *testing.T) {
	data := map[string]string{
		COLUMN_ID:              "test-id",
		COLUMN_ENTITY_TYPE:     "test-entity-type",
		COLUMN_ENTITY_ID:       "test-entity-id",
		COLUMN_CONTENT:         "test-content",
		COLUMN_CREATED_AT:      "2024-01-01 00:00:00",
		COLUMN_SOFT_DELETED_AT: sb.MAX_DATETIME,
	}

	version := NewVersionFromExistingData(data)
	result := version.Data()

	if result[COLUMN_ID] != "test-id" {
		t.Errorf("Data()[ID] = %s, want %s", result[COLUMN_ID], "test-id")
	}

	if result[COLUMN_ENTITY_TYPE] != "test-entity-type" {
		t.Errorf("Data()[EntityType] = %s, want %s", result[COLUMN_ENTITY_TYPE], "test-entity-type")
	}

	if result[COLUMN_ENTITY_ID] != "test-entity-id" {
		t.Errorf("Data()[EntityID] = %s, want %s", result[COLUMN_ENTITY_ID], "test-entity-id")
	}

	if result[COLUMN_CONTENT] != "test-content" {
		t.Errorf("Data()[Content] = %s, want %s", result[COLUMN_CONTENT], "test-content")
	}

	if result[COLUMN_CREATED_AT] != "2024-01-01 00:00:00" {
		t.Errorf("Data()[CreatedAt] = %s, want %s", result[COLUMN_CREATED_AT], "2024-01-01 00:00:00")
	}

	if result[COLUMN_SOFT_DELETED_AT] != sb.MAX_DATETIME {
		t.Errorf("Data()[SoftDeletedAt] = %s, want %s", result[COLUMN_SOFT_DELETED_AT], sb.MAX_DATETIME)
	}
}

func TestVersionDataChanged(t *testing.T) {
	version := NewVersion().
		SetSoftDeletedAt("2024-12-31 23:59:59")

	dataChanged := version.DataChanged()

	if dataChanged[COLUMN_SOFT_DELETED_AT] != "2024-12-31 23:59:59" {
		t.Errorf("DataChanged()[SoftDeletedAt] = %s, want %s", dataChanged[COLUMN_SOFT_DELETED_AT], "2024-12-31 23:59:59")
	}

	if len(dataChanged) != 1 {
		t.Errorf("DataChanged() should return 1 field, got %d", len(dataChanged))
	}
}

func TestVersionMarkAsNotDirty(t *testing.T) {
	version := NewVersion()

	// Should not panic or error
	version.MarkAsNotDirty()
}

func TestVersionContent(t *testing.T) {
	version := NewVersion()

	if version.Content() != "" {
		t.Errorf("Content() should be empty initially, got %s", version.Content())
	}

	result := version.SetContent("new-content")

	if result != version {
		t.Error("SetContent() should return the same instance for chaining")
	}

	if version.Content() != "new-content" {
		t.Errorf("Content() = %s, want %s", version.Content(), "new-content")
	}
}

func TestVersionCreatedAt(t *testing.T) {
	version := NewVersion()

	initialCreatedAt := version.CreatedAt()
	if initialCreatedAt == "" {
		t.Error("CreatedAt() should be set by NewVersion()")
	}

	result := version.SetCreatedAt("2024-06-15 12:30:45")

	if result != version {
		t.Error("SetCreatedAt() should return the same instance for chaining")
	}

	if version.CreatedAt() != "2024-06-15 12:30:45" {
		t.Errorf("CreatedAt() = %s, want %s", version.CreatedAt(), "2024-06-15 12:30:45")
	}
}

func TestVersionEntityType(t *testing.T) {
	version := NewVersion()

	if version.EntityType() != "" {
		t.Errorf("EntityType() should be empty initially, got %s", version.EntityType())
	}

	result := version.SetEntityType("page")

	if result != version {
		t.Error("SetEntityType() should return the same instance for chaining")
	}

	if version.EntityType() != "page" {
		t.Errorf("EntityType() = %s, want %s", version.EntityType(), "page")
	}
}

func TestVersionEntityID(t *testing.T) {
	version := NewVersion()

	if version.EntityID() != "" {
		t.Errorf("EntityID() should be empty initially, got %s", version.EntityID())
	}

	result := version.SetEntityID("12345")

	if result != version {
		t.Error("SetEntityID() should return the same instance for chaining")
	}

	if version.EntityID() != "12345" {
		t.Errorf("EntityID() = %s, want %s", version.EntityID(), "12345")
	}
}

func TestVersionID(t *testing.T) {
	version := NewVersion()

	initialID := version.ID()
	if initialID == "" {
		t.Error("ID() should be set by NewVersion()")
	}

	result := version.SetID("custom-id")

	if result != version {
		t.Error("SetID() should return the same instance for chaining")
	}

	if version.ID() != "custom-id" {
		t.Errorf("ID() = %s, want %s", version.ID(), "custom-id")
	}
}

func TestVersionSoftDeletedAt(t *testing.T) {
	version := NewVersion()

	if version.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Errorf("SoftDeletedAt() should be MAX_DATETIME initially, got %s", version.SoftDeletedAt())
	}

	result := version.SetSoftDeletedAt("2024-12-31 23:59:59")

	if result != version {
		t.Error("SetSoftDeletedAt() should return the same instance for chaining")
	}

	if version.SoftDeletedAt() != "2024-12-31 23:59:59" {
		t.Errorf("SoftDeletedAt() = %s, want %s", version.SoftDeletedAt(), "2024-12-31 23:59:59")
	}
}

func TestVersionMethodChaining(t *testing.T) {
	version := NewVersion().
		SetID("chain-id").
		SetEntityType("chain-type").
		SetEntityID("chain-entity-id").
		SetContent("chain-content").
		SetCreatedAt("2024-01-01 00:00:00").
		SetSoftDeletedAt("2024-12-31 23:59:59")

	if version.ID() != "chain-id" {
		t.Errorf("Method chaining failed for ID, got %s", version.ID())
	}

	if version.EntityType() != "chain-type" {
		t.Errorf("Method chaining failed for EntityType, got %s", version.EntityType())
	}

	if version.EntityID() != "chain-entity-id" {
		t.Errorf("Method chaining failed for EntityID, got %s", version.EntityID())
	}

	if version.Content() != "chain-content" {
		t.Errorf("Method chaining failed for Content, got %s", version.Content())
	}

	if version.CreatedAt() != "2024-01-01 00:00:00" {
		t.Errorf("Method chaining failed for CreatedAt, got %s", version.CreatedAt())
	}

	if version.SoftDeletedAt() != "2024-12-31 23:59:59" {
		t.Errorf("Method chaining failed for SoftDeletedAt, got %s", version.SoftDeletedAt())
	}
}
