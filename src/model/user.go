package model

type User struct {
	Id        int64  `json:"id" db:"id" goqu:"skipinsert"`
	Login     string `json:"login" db:"login"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}
