package models

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserAccount struct {
	Email       string `json:"email"  binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type UserResponse struct {
	LastName    string `json:"lastname,omitempty"`
	FirstName   string `json:"firstname,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phonenumber,omitempty"`
}

type AuthToken struct {
	Token string `json:"token,omitempty"`
}
