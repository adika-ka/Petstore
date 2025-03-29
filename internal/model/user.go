package model

type User struct {
	ID         int64  `db:"id" json:"id" example:"1"`
	Username   string `db:"username" json:"username" example:"johndoe"`
	FirstName  string `db:"first_name" json:"firstName" example:"John"`
	LastName   string `db:"last_name" json:"lastName" example:"Doe"`
	Email      string `db:"email" json:"email" example:"johndoe@example.com"`
	Password   string `db:"password" json:"password" example:"secret123"`
	Phone      string `db:"phone" json:"phone" example:"+123456789"`
	UserStatus int    `db:"user_status" json:"userStatus" example:"1"`
}
