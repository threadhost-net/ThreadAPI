package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserSession struct {
	UserAgent string
	IP        string
	CreatedAt time.Time
}

type SessionInfo struct {
	userAgent string
	ip        string
	createdAt time.Time
	isCurrent bool
}

func IssueToken(ctx context.Context, client *redis.Client, userID, value string, ttl time.Duration) (string, error) {
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

func RevokeAllForUser(ctx context.Context, client *redis.Client, userID int64) error {
	iter := client.Scan(ctx, 0, "session:*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		parts := strings.SplitN(key, ":", 3)
		if len(parts) < 3 {
			continue
		}

		if strconv.FormatInt(userID, 10) == parts[1] {
			if err := client.Del(ctx, key).Err(); err != nil {
				log.Println("Del session failed: ", err)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func List(ctx context.Context, client *redis.Client, userID int64, currentToken string) ([]SessionInfo, error) {
	var sessions []SessionInfo
	iter := client.Scan(ctx, 0, "session:*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		raw, err := client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var s UserSession
		if err := json.Unmarshal([]byte(raw), &s); err != nil {
			continue
		}

		parts := strings.SplitN(key, ":", 3)
		if len(parts) < 3 {
			continue
		}

		if parts[1] == strconv.FormatInt(userID, 10) {
			isCurrent := "session:"+parts[1]+":"+currentToken == key
			sessions = append(sessions, SessionInfo{
				userAgent: s.UserAgent,
				ip:        s.IP,
				createdAt: s.CreatedAt,
				isCurrent: isCurrent,
			})
		}
	}
	if err := iter.Err(); err != nil {
		return sessions, err
	}
	return sessions, nil
}
