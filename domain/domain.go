package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/clerk/clerk-sdk-go/v2"
)

type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

func (m *BaseModel) IsValid() bool {
	if m.ID == "" {
		return false
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		return false
	}
	return true
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New().String()
	m.CreatedAt = time.Now().UTC().Local()
	m.UpdatedAt = time.Now().UTC().Local()
	if !m.IsValid() {
		return errors.New("rollback: invalid model")
	}
	return
}

type Household struct {
	BaseModel
	Name string `json:"name"`
}

type Member struct {
	BaseModel
	UserID      string `json:"userId"`
	HouseholdID string `json:"householdId"`
}

type Meal struct {
	BaseModel
	Name        string `json:"name"`
	HouseholdID string `json:"householdId"`
}

type Event struct {
	BaseModel
	Name       string    `json:"name"`
	EntityID   string    `json:"entityId"`   // Could be HouseholdID or MealID etc.
	EntityType string    `json:"entityType"` // e.g., "household", "meal", "chore"
	StartDate  time.Time `json:"startDate" gorm:"autoUpdateTime:false"`
	EndDate    time.Time `json:"endDate" gorm:"autoUpdateTime:false"`
	AssignedTo string    `json:"assignedTo"` // UserID of the person assigned to this event
}

type User struct {
	BaseModel
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	ID                  string              `json:"id"`
	Username            string              `json:"username"`
	FirstName           string              `json:"firstName"`
	LastName            string              `json:"lastName"`
	ProfileImageURL     string              `json:"profileImageUrl"`
	PrimaryEmailAddress string              `json:"primaryEmailAddress"`
	EmailAddresses      *clerk.EmailAddress `json:"emailAddresses"`
}
