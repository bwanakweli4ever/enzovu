package models

// User represents the model for the User resource.
type User struct {
	// ID is the primary identifier for the resource
	ID    int    `json:"id"`

	// Name is the name of the resource
	Name  string `json:"name"`

	// Email is the email of the resource owner (optional)
	Email string `json:"email"`

	// Add more fields as necessary (e.g., Description, CreatedAt, UpdatedAt)
}

// GetUser returns a sample User instance
func GetUser() *User {
	return &User{
		ID:    1,
		Name:  "Example %!s(MISSING)",
		Email: "example@email.com",
	}
}
