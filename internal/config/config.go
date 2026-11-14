package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string

	RedisHost     string
	RedisPassword string

	ServerAddress string
	Domain        string

	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	SmtpFrom     string

	IsLeader bool

	HCaptchaSecret  string
	HCaptchaSitekey string

	PaymentsAPIAddress string
}

func New() *Config {
	return &Config{
		DatabaseURL:   envOrDefault("DATABASE_URL", "postgres://user:password@localhost:5432/database?sslmode=disable"),
		RedisHost:     envOrDefault("REDIS_HOST", "localhost:6379"),
		RedisPassword: envOrDefault("REDIS_PASSWORD", "password"),

		ServerAddress: envOrDefault("ADDRESS", "0.0.0.0:3000"),
		Domain:        envOrDefault("DOMAIN", "localhost"),

		SmtpHost:     envOrDefault("SMTP_HOST", "smtp.example.com"),
		SmtpPort:     envOrDefaultInt("SMTP_PORT", 587),
		SmtpUsername: envOrDefault("SMTP_USERNAME", "user"),
		SmtpPassword: envOrDefault("SMTP_PASSWORD", "password"),
		SmtpFrom:     envOrDefault("SMTP_FROM", "user@example.com"),

		IsLeader: envOrDefault("IS_LEADER", "true") == "true",

		HCaptchaSecret:  envOrDefault("HCAPTCHA_SECRET", ""),
		HCaptchaSitekey: envOrDefault("HCAPTCHA_SITEKEY", ""),

		PaymentsAPIAddress: envOrDefault("PAYMETSAPI_ADDRESS", "localhost"),
	}
}

func envOrDefault(key string, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func envOrDefaultInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
