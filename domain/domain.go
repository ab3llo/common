package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string         `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
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
