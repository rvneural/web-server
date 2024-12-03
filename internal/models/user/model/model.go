package model

type User struct {
	Email     string
	FirstName string
	LastName  string
	ID        string
}

type AuthentificationInput struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegistrationInput struct {
	Email     string `json:"email" form:"email" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required"`
	FirstName string `json:"firstName" form:"firstName" binding:"required"`
	LastName  string `json:"lastName" form:"lastName" binding:"required"`
}
