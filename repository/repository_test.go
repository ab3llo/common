package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&TestModel{})
	return db
}

func TestRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository[TestModel](db)
	ctx := context.Background()

	// Create
	model := &TestModel{ID: "1", Name: "test"}
	created, err := repo.Create(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, model, created)

	// Get
	got, err := repo.Get(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, model.ID, got.ID)

	// GetAll
	all, err := repo.GetAll(ctx, 10, 0)
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	// Update
	model.Name = "updated"
	updated, err := repo.Update(ctx, "1", model)
	assert.NoError(t, err)
	assert.Equal(t, "updated", updated.Name)

	// Delete
	err = repo.Delete(ctx, "1")
	assert.NoError(t, err)

	// Get after delete
	_, err = repo.Get(ctx, "1")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepository_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository[TestModel](db)
	ctx := context.Background()
	_, err := repo.Get(ctx, "not-exist")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
