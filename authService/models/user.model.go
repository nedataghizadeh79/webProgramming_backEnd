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
	UserId      string `json:"userId,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type AuthToken struct {
	Token string `json:"token,omitempty"`
}
