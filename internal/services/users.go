package services

import (
	"errors"
	"sso/internal/models"
	"sso/internal/storage"
	"sso/pkg/helpers/errorHelper"
	"sso/pkg/helpers/passwdHelper"
	"strings"
	"unicode"
)

const (
	ErrorCreatePassword   = "Password hashing error"
	ErrorPasswordValidate = "Password validation error"
	ErrorAddUser          = "Error on add user"
)

type UsersService struct {
	storage storage.Storage
}

func Users(storage storage.Storage) *UsersService {
	return &UsersService{
		storage: storage,
	}
}

func (u *UsersService) Add(user *models.User) (string, error) {
	const operation = "internal.services.users.Add()"
	if err := validatePassword(user.Password); err != nil {
		return "", errorHelper.WrapError(operation, ErrorPasswordValidate, err)
	} else {
		password, err := passwdHelper.HashPassword(user.Password)
		if err != nil {
			return "", errorHelper.WrapError(operation, ErrorCreatePassword, err)
		}
		user.Password = password
		uid, err := u.storage.Users().InsertUser(user)
		if err != nil {
			return "", errorHelper.WrapError(operation, ErrorAddUser, err)
		}
		return uid, nil
	}
}

// Password validates plain password against the rules defined below.
//
// TODO переписать на нормальные валидаторы
// upp: at least one upper case letter.
// low: at least one lower case letter.
// num: at least one digit.
// sym: at least one special character.
// tot: at least eight characters long.
// No empty string or whitespace.
func validatePassword(pass string) error {
	var (
		upp, low, num, sym bool
		tot                uint8
	)
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return errors.New("Empty password")
		}
	}
	if !upp || !low || !num || !sym || tot < 8 {
		b := strings.Builder{}
		if !upp {
			b.WriteString("Password must contain one or more uppercase letter. ")
		}
		if !low {
			b.WriteString("Password must contain one or more lowercase letter. ")
		}
		if !num {
			b.WriteString("The password must contain one or more digits. ")
		}
		if !sym {
			b.WriteString("The password must contain one or more special character. ")
		}
		if tot < 8 {
			b.WriteString("The password must be at least 8 characters long. ")
		}
		return errors.New(b.String()[:b.Len()-1])
	}
	return nil
}
