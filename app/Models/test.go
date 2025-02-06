package models

type Test struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Example function to return a Test instance
func GetTest() *Test {
	return &Test{ID: 1, Name: "Example %!s(MISSING)", Email: "example@email.com"}
}
