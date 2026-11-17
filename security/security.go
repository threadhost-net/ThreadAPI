package security

import (
	"errors"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]+$`)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("Email обязателен")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("Некорректный email")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("Имя пользователя слишком короткое")
	}
	if len(username) > 32 {
		return errors.New("Имя пользователя слишком длинное")
	}
	if !usernameRe.MatchString(username) {
		return errors.New("Имя пользователя содержит запрещённые символы")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Пароль слишком короткий")
	}
	if len(password) > 128 {
		return errors.New("Пароль слишком длинный")
	}
	return nil
}

func Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
