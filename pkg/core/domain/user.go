package domain

type User struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	DateCreated  string `json:"date_created"`
	LastModified string `json:"-"`
	Status       string `json:"-"`
	Role         string `json:"role"`
}
