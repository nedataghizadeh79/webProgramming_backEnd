package utils

import (
	"AuthService/models"
	"errors"
	"regexp"
)

var (
	persianRegex  = `^[\u0600-\u06FF]+[\s\u0600-\u06FF]*$`
	emailRegex    = `^([\w\.\_]{2,10})@(\w{1,})\.([a-z]{2,4})$`
	passwordRegex = `^\w{6,}$`
	phoneRegex    = `^(\+98|09)\d{9}$`
)

func ValidateSignUpData(user models.UserAccount) error {

	if len(user.Email) == 0 {
		return errors.New("your email should not be empty")
	}

	if len(user.Password) == 0 {
		return errors.New("your password should not be empty")
	}

	if len(user.Gender) == 0 {
		return errors.New("your gender should not be empty")
	}

	if len(user.PhoneNumber) == 0 {
		return errors.New("your phone number should not be empty")
	}

	if len(user.FirstName) == 0 {
		return errors.New("your first name should not be empty")
	}

	if len(user.LastName) == 0 {
		return errors.New("your last name should not be empty")
	}

	if m, _ := regexp.MatchString(persianRegex, user.FirstName); !m {
		return errors.New("invalid name")
	}

	if m, _ := regexp.MatchString(persianRegex, user.LastName); !m {
		return errors.New("invalid last name")
	}

	if m, _ := regexp.MatchString("^(F|M)$", user.Gender); !m {
		return errors.New("your gender should be F or M")
	}

	if m, _ := regexp.MatchString(phoneRegex, user.PhoneNumber); !m {
		return errors.New("invalid phone number")
	}

	if m, _ := regexp.MatchString(emailRegex, user.Email); !m {
		return errors.New("your email is not valid")
	}

	if m, _ := regexp.MatchString(passwordRegex, user.Password); !m {
		return errors.New("invalid password")
	}

	return nil
}

func ValidateSignInData(user models.SignInInput) error {

	if len(user.Email) == 0 {
		return errors.New("your email should not be empty")
	}

	if len(user.Password) == 0 {
		return errors.New("your password should not be empty")
	}

	if m, _ := regexp.MatchString(emailRegex, user.Email); !m {
		return errors.New("your email is not valid")
	}

	if m, _ := regexp.MatchString(passwordRegex, user.Password); !m {
		return errors.New("invalid password")
	}

	return nil
}
