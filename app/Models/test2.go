package models

// Test2 represents the model for the Test2 resource.
type Test2 struct {
	// ID is the primary identifier for the resource
	ID    int    `json:"id"`

	// Name is the name of the resource
	Name  string `json:"name"`

	// Email is the email of the resource owner (optional)
	Email string `json:"email"`

	// Add more fields as necessary (e.g., Description, CreatedAt, UpdatedAt)
}

// GetTest2 returns a sample Test2 instance
func GetTest2() *Test2 {
	return &Test2{
		ID:    1,
		Name:  "Example %!s(MISSING)",
		Email: "example@email.com",
	}
}
