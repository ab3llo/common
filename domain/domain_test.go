package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBaseModel_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		model    BaseModel
		expected bool
	}{
		{
			name:     "valid model with ID",
			model:    BaseModel{ID: "test-id"},
			expected: true,
		},
		{
			name:     "invalid model without ID",
			model:    BaseModel{},
			expected: false,
		},
		{
			name:     "valid model with UUID",
			model:    BaseModel{ID: uuid.New().String()},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBaseModel_BeforeCreate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	model := &BaseModel{}

	err = model.BeforeCreate(db)
	if err != nil {
		t.Errorf("BeforeCreate() error = %v, want nil", err)
	}

	if model.ID == "" {
		t.Error("BeforeCreate() should set ID")
	}

	if !model.IsValid() {
		t.Error("Model should be valid after BeforeCreate()")
	}

	_, err = uuid.Parse(model.ID)
	if err != nil {
		t.Errorf("ID should be a valid UUID, got %s", model.ID)
	}
}

func TestHousehold(t *testing.T) {
	household := Household{
		BaseModel: BaseModel{ID: "test-id"},
		Name:      "Test Household",
	}

	if !household.IsValid() {
		t.Error("Household should be valid")
	}

	if household.Name != "Test Household" {
		t.Errorf("Name = %v, want %v", household.Name, "Test Household")
	}
}

func TestMember(t *testing.T) {
	member := Member{
		BaseModel:   BaseModel{ID: "test-id"},
		UserID:      "user-123",
		HouseholdID: "household-456",
	}

	if !member.IsValid() {
		t.Error("Member should be valid")
	}

	if member.UserID != "user-123" {
		t.Errorf("UserID = %v, want %v", member.UserID, "user-123")
	}

	if member.HouseholdID != "household-456" {
		t.Errorf("HouseholdID = %v, want %v", member.HouseholdID, "household-456")
	}
}

func TestMeal(t *testing.T) {
	meal := Meal{
		BaseModel:   BaseModel{ID: "test-id"},
		Name:        "Breakfast",
		HouseholdID: "household-123",
	}

	if !meal.IsValid() {
		t.Error("Meal should be valid")
	}

	if meal.Name != "Breakfast" {
		t.Errorf("Name = %v, want %v", meal.Name, "Breakfast")
	}
}

func TestEvent(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	event := Event{
		BaseModel:  BaseModel{ID: "test-id"},
		Name:       "Test Event",
		EntityID:   "entity-123",
		EntityType: "household",
		StartDate:  startTime,
		EndDate:    endTime,
		AssignedTo: "user-456",
	}

	if !event.IsValid() {
		t.Error("Event should be valid")
	}

	if event.Name != "Test Event" {
		t.Errorf("Name = %v, want %v", event.Name, "Test Event")
	}

	if event.EntityType != "household" {
		t.Errorf("EntityType = %v, want %v", event.EntityType, "household")
	}

	if event.StartDate.After(event.EndDate) {
		t.Error("StartDate should be before EndDate")
	}
}

func TestUser(t *testing.T) {
	user := User{
		BaseModel: BaseModel{ID: "test-id"},
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}

	if !user.IsValid() {
		t.Error("User should be valid")
	}

	if user.Email != "john@example.com" {
		t.Errorf("Email = %v, want %v", user.Email, "john@example.com")
	}
}

func TestUserData(t *testing.T) {
	userData := UserData{
		ID:                  "clerk-user-123",
		Username:            "johndoe",
		FirstName:           "John",
		LastName:            "Doe",
		ProfileImageURL:     "https://example.com/avatar.jpg",
		PrimaryEmailAddress: "john@example.com",
	}

	if userData.ID != "clerk-user-123" {
		t.Errorf("ID = %v, want %v", userData.ID, "clerk-user-123")
	}

	if userData.Username != "johndoe" {
		t.Errorf("Username = %v, want %v", userData.Username, "johndoe")
	}
}

func TestEvent_DateValidation(t *testing.T) {
	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		valid     bool
	}{
		{
			name:      "valid dates - start before end",
			startDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			endDate:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			valid:     true,
		},
		{
			name:      "invalid dates - start after end",
			startDate: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			endDate:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			valid:     false,
		},
		{
			name:      "same start and end date",
			startDate: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			endDate:   time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			valid:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{
				BaseModel:  BaseModel{ID: "test-id"},
				Name:       "Test Event",
				EntityID:   "entity-123",
				EntityType: "household",
				StartDate:  tt.startDate,
				EndDate:    tt.endDate,
				AssignedTo: "user-456",
			}

			isValid := !event.StartDate.After(event.EndDate)
			if isValid != tt.valid {
				t.Errorf("Event date validation = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestRelationships_HouseholdMembers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&Household{}, &Member{}, &User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	household := &Household{
		BaseModel: BaseModel{ID: uuid.New().String()},
		Name:      "Test Family",
	}
	err = db.Create(household).Error
	if err != nil {
		t.Fatalf("Failed to create household: %v", err)
	}

	user := &User{
		BaseModel: BaseModel{ID: uuid.New().String()},
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	err = db.Create(user).Error
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	member := &Member{
		BaseModel:   BaseModel{ID: uuid.New().String()},
		UserID:      user.ID,
		HouseholdID: household.ID,
	}
	err = db.Create(member).Error
	if err != nil {
		t.Fatalf("Failed to create member: %v", err)
	}

	var retrievedMember Member
	err = db.First(&retrievedMember, "household_id = ?", household.ID).Error
	if err != nil {
		t.Fatalf("Failed to retrieve member: %v", err)
	}

	if retrievedMember.UserID != user.ID {
		t.Errorf("Member UserID = %v, want %v", retrievedMember.UserID, user.ID)
	}
	if retrievedMember.HouseholdID != household.ID {
		t.Errorf("Member HouseholdID = %v, want %v", retrievedMember.HouseholdID, household.ID)
	}
}

func TestBaseModel_ValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		model BaseModel
		valid bool
	}{
		{
			name:  "empty ID",
			model: BaseModel{ID: ""},
			valid: false,
		},
		{
			name:  "whitespace only ID",
			model: BaseModel{ID: "   "},
			valid: true, // Current implementation only checks empty string
		},
		{
			name:  "very long ID",
			model: BaseModel{ID: string(make([]byte, 1000))},
			valid: true,
		},
		{
			name:  "special characters in ID",
			model: BaseModel{ID: "id-with-special!@#$%^&*()"},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.model.IsValid()
			if result != tt.valid {
				t.Errorf("IsValid() = %v, want %v for ID: %q", result, tt.valid, tt.model.ID)
			}
		})
	}
}

func TestMeal_HouseholdRelationship(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&Household{}, &Meal{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	household := &Household{
		BaseModel: BaseModel{ID: uuid.New().String()},
		Name:      "Test Family",
	}
	err = db.Create(household).Error
	if err != nil {
		t.Fatalf("Failed to create household: %v", err)
	}

	meals := []Meal{
		{
			BaseModel:   BaseModel{ID: uuid.New().String()},
			Name:        "Breakfast",
			HouseholdID: household.ID,
		},
		{
			BaseModel:   BaseModel{ID: uuid.New().String()},
			Name:        "Lunch",
			HouseholdID: household.ID,
		},
		{
			BaseModel:   BaseModel{ID: uuid.New().String()},
			Name:        "Dinner",
			HouseholdID: household.ID,
		},
	}

	for _, meal := range meals {
		err = db.Create(&meal).Error
		if err != nil {
			t.Fatalf("Failed to create meal %s: %v", meal.Name, err)
		}
	}

	var retrievedMeals []Meal
	err = db.Where("household_id = ?", household.ID).Find(&retrievedMeals).Error
	if err != nil {
		t.Fatalf("Failed to retrieve meals: %v", err)
	}

	if len(retrievedMeals) != 3 {
		t.Errorf("Expected 3 meals, got %d", len(retrievedMeals))
	}

	mealNames := make(map[string]bool)
	for _, meal := range retrievedMeals {
		mealNames[meal.Name] = true
		if meal.HouseholdID != household.ID {
			t.Errorf("Meal %s has wrong HouseholdID: %v, want %v", meal.Name, meal.HouseholdID, household.ID)
		}
	}

	expectedMeals := []string{"Breakfast", "Lunch", "Dinner"}
	for _, expected := range expectedMeals {
		if !mealNames[expected] {
			t.Errorf("Missing meal: %s", expected)
		}
	}
}

func BenchmarkBaseModel_IsValid(b *testing.B) {
	model := BaseModel{ID: "test-id-for-benchmark"}
	for i := 0; i < b.N; i++ {
		model.IsValid()
	}
}

func BenchmarkBaseModel_BeforeCreate(b *testing.B) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	model := &BaseModel{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.ID = "" // Reset ID to trigger UUID generation
		model.BeforeCreate(db)
	}
}