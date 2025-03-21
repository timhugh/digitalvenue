package dv

type User struct {
	UUID         string `json:"uuid" db:"uuid"`
	Name         string `json:"name" db:"name"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}
