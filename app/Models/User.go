package models

// User struct defines the structure of a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUser returns a mock user
func GetUser() *User {
	return &User{
		ID:    1,
		Name:  "John Doe",
		Email: "johndoe@example.com",
	}
}
