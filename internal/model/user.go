package model

type User struct {
	ID         int64  `db:"id" json:"id"`
	Username   string `db:"username" json:"username"`
	FirstName  string `db:"first_name" json:"firstName"`
	LastName   string `db:"last_name" json:"lastName"`
	Email      string `db:"email" json:"email"`
	Password   string `db:"password" json:"password"`
	Phone      string `db:"phone" json:"phone"`
	UserStatus int    `db:"user_status" json:"userStatus"`
}
