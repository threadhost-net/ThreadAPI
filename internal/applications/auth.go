package applications

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/threadhost-net/ThreadAPI/assets/emails"
	"github.com/threadhost-net/ThreadAPI/security"
	"github.com/threadhost-net/ThreadAPI/session"
)

func (tb *Toolbox) Register(ctx context.Context, username, password, email, hCaptchaResponse string) (status int, err error) {
	// TODO капча

	if err := security.ValidateUsername(username); err != nil {
		return http.StatusBadRequest, err
	}

	if err := security.ValidatePassword(password); err != nil {
		return http.StatusBadRequest, err
	}

	if err := security.ValidateEmail(email); err != nil {
		return http.StatusBadRequest, err
	}

	hashPassword, err := security.Hash(password)
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

	emailToken, err := session.IssueToken(ctx, tb.Redis.Get("email_verify"), "token:", strconv.FormatInt(userID, 10), 24*time.Hour)
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

func (tb *Toolbox) Login(ctx context.Context, username, password, hCaptchaResponse, address, userAgent string) (status int, cookieToken string, err error) {
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

	if !security.VerifyPassword(passwordHash, password) {
		return http.StatusUnauthorized, "", errors.New("Неверный пароль")
	}

	if !isActivated {
		return http.StatusForbidden, "", errors.New("Учётная запись не подтверждена")
	}

	currentSession := session.UserSession{
		UserAgent: userAgent,
		IP:        address,
		CreatedAt: time.Now().UTC(),
	}

	jsonSession, err := json.Marshal(currentSession)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, "", errors.New("Не удалось создать сессию")
	}

	userIDStr := strconv.FormatInt(userID, 10)
	token, err := session.IssueToken(ctx, tb.Redis.Get("sessions"), userIDStr, string(jsonSession), 7*24*time.Hour)
	if err != nil {
		return http.StatusInternalServerError, "", errors.New("Не удалось создать сессию")
	}

	cookieToken = userIDStr + ":" + token
	return http.StatusAccepted, cookieToken, nil
}

func (tb *Toolbox) Logout(ctx context.Context, cookieToken string) (status int, err error) {
	client := tb.Redis.Get("sessions")
	if err := client.Del(ctx, "session:"+cookieToken).Err(); err != nil {
		log.Println(err)
		return http.StatusInternalServerError, errors.New("Не удалось удалить вашу сессию с сервера")
	}
	return http.StatusOK, nil
}
