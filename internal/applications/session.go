package applications

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type session struct {
	userAgent string
	ip        string
	createdAt time.Time
}

type sessionInfo struct {
	userAgent string
	ip        string
	createdAt time.Time
	isCurrent bool
}

func revokeAllForUser(ctx context.Context, client *redis.Client, userID int64) error {
	iter := client.Scan(ctx, 0, "session:*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		raw, err := client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var s session
		if err := json.Unmarshal([]byte(raw), &s); err != nil {
			continue
		}
		if s.UserID == userID {
			if err := client.Del(ctx, key).Err(); err != nil {
				log.Println("del session: ", err)
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func List(ctx context.Context, client *redis.Client, userID int64, currentToken string) ([]sessionInfo, error) {
	var sessions []sessionInfo
	iter := client.Scan(ctx, 0, "session:*", 100).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		raw, err := client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var s session
		if err := json.Unmarshal([]byte(raw), &s); err != nil {
			continue
		}
		if s.UserID == userID {
			isCurrent := "session:"+currentToken == key
			sessions = append(sessions, sessionInfo{
				UserAgent: s.UserAgent,
				IP:        s.IP,
				CreatedAt: s.CreatedAt,
				IsCurrent: isCurrent,
			})
		}
	}
	if err := iter.Err(); err != nil {
		return sessions, err
	}
	return sessions, nil
}
