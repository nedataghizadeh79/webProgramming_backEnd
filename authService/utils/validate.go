package utils

import (
	"AuthService/models"
	"errors"
	"regexp"
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

	//should be english
	if m, _ := regexp.MatchString("^[a-zA-Z]+$", user.FirstName); !m {
		return errors.New("your language should be English ")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z]+$", user.LastName); !m {
		return errors.New("your language should be English ")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z]+$", user.Password); !m {
		return errors.New("your language should be English ")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z]+$", user.Gender); !m {
		return errors.New("your language should be English ")
	}

	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, user.Email); !m {
		return errors.New("your email is not valid")
	}

	if m, _ := regexp.MatchString(`^[_A-Za-z0-9+-]+([.][_A-Za-z0-9-]+)*@[A-Za-z0-9-]+([.][A-Za-z0-9]+)*([.][A-Za-z]{2,})$`, user.Password); !m {
		return errors.New("invalid password")
	}

	return nil
}
