package domain

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at,omitempty"`
}
type Household struct {
	BaseModel
	Name string `json:"name"`
}

func (m *BaseModel) IsValid() bool {
	if m.ID == "" {
		return false
	}
	return true
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New().String()

	if !m.IsValid() {
		return errors.New("rollback: invalid model")
	}
	return
}
