package applications

import (
	"errors"
	"net/mail"
	"regexp"
)

// Паттерн для проверки имени пользователя
var usernameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]+$`)

func validateEmail(email string) error {
	if email == "" {
		return errors.New("Email обязателен")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("Некорректный email")
	}
	return nil
}

func validateUsername(username string) error {
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

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Пароль слишком короткий")
	}
	if len(password) > 128 {
		return errors.New("Пароль слишком длинный")
	}
	return nil
}
