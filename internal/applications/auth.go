package applications

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"github.com/threadhost-net/ThreadAPI/assets/emails"
	"golang.org/x/crypto/bcrypt"
)

func (tb *Toolbox) Register(ctx context.Context, username, password, email, hCaptchaResponse string) (code int, err error) {
	// TODO капча

	if err := validateUsername(username); err != nil {
		return http.StatusBadRequest, err
	}

	if err := validatePassword(password); err != nil {
		return http.StatusBadRequest, err
	}

	if err := validateEmail(email); err != nil {
		return http.StatusBadRequest, err
	}

	hashPassword, err := hash(password)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Не удалось захэшировать пароль")
	}

	var userID int64
	if err := tb.DB.QueryRow(ctx, "INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id", username, hashPassword, email).Scan(&userID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return http.StatusConflict, errors.New("Имя пользователя или почта уже заняты")
		}
		log.Println(err)
		return http.StatusInternalServerError, errors.New("Не удалось создать учётную запись")
	}

	emailToken, err := issueToken(ctx, tb.Redis.Get("email_verify"), "token:", strconv.FormatInt(userID, 10), 24*time.Hour)
	if err != nil {
		log.Println(err)

	} else {
		body, err := emails.Render("email_confirm_ru.html", map[string]any{
			"Token":    emailToken,
			"Username": username,
		})
		if err != nil {
			log.Println("render email failed: ", err)
		} else if err := tb.Mailer.Send(email, "Подтверждение почты", body); err != nil {
			log.Println("send email failed: ", err)
		}
	}
	return http.StatusCreated, nil
}

func (tb *Toolbox) Login(ctx context.Context, username, password, hCaptchaResponse, address, userAgent string) (code int, cookieToken string, err error) {
	// TODO капча
	var (
		userID       int64
		passwordHash string
		isActivated  bool
	)

	err = tb.DB.QueryRow(ctx, "SELECT id, password_hash, is_activated FROM users WHERE username = $1", username).Scan(&userID, &passwordHash, &isActivated)
	if errors.Is(err, pgx.ErrNoRows) {
		return http.StatusUnauthorized, "", errors.New("Учётной записи не существует")
	}
	if err != nil {
		log.Println("err")
		return http.StatusInternalServerError, "", errors.New("Не удалось проверить вашу учётную запись")
	}

	if !verifyPassword(passwordHash, password) {
		return http.StatusUnauthorized, "", errors.New("Неверный пароль")
	}

	if !isActivated {
		return http.StatusForbidden, "", errors.New("Учётная запись не подтверждена")
	}

	currentSession := session{
		userAgent: userAgent,
		ip:        address,
		createdAt: time.Now().UTC(),
	}

	jsonSession, err := json.Marshal(currentSession)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "", errors.New("Не удалось создать сессию")
	}

	userIDStr := strconv.FormatInt(userID, 10)
	token, err := issueToken(ctx, tb.Redis.Get("sessions"), userIDStr, string(jsonSession), 7*24*time.Hour)
	if err != nil {
		return http.StatusInternalServerError, "", errors.New("Не удалось создать сессию")
	}

	return http.StatusAccepted, token, nil
}

// Утилиты

func hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func verifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func issueToken(ctx context.Context, client *redis.Client, userID, value string, ttl time.Duration) (string, error) {
	for {
		token, err := newToken()
		if err != nil {
			return "", err
		}

		ok, err := client.SetNX(ctx, userID+":"+token, value, ttl).Result()
		if err != nil {
			log.Println(err)
			return "", err
		}
		if ok {
			return token, nil
		}
	}
}

func newToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
