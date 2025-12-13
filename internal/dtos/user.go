package dtos

type Password struct {
	Text *string
	Hash []byte
}

type InsertUser struct {
	ID           string   `json:"-"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Email        string   `json:"email"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	PasswordHash Password `json:"-"`
	CreatedAt    int64    `json:"-"`
	UpdatedAt    int64    `json:"-"`
	Version      int      `json:"-"`
}
