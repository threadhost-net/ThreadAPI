package config

import "github.com/threadhost-net/ThreadAPI/env"

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
		DatabaseURL: env.OrDefault("DATABASE_URL", "postgres://user:password@localhost:5432/database?sslmode=disable"),

		RedisHost:     env.OrDefault("REDIS_HOST", "localhost:6379"),
		RedisPassword: env.OrDefault("REDIS_PASSWORD", "password"),

		ServerAddress: env.OrDefault("ADDRESS", "0.0.0.0:3000"),
		Domain:        env.OrDefault("DOMAIN", "localhost"),

		SmtpHost:     env.OrDefault("SMTP_HOST", "smtp.example.net"),
		SmtpPort:     env.OrDefaultInt("SMTP_PORT", 465),
		SmtpUsername: env.OrDefault("SMTP_USERNAME", "mail@example.net"),
		SmtpPassword: env.OrDefault("SMTP_PASSWORD", "password"),
		SmtpFrom:     env.OrDefault("SMTP_FROM", "noreply@example.net"),

		IsLeader: env.OrDefault("IS_LEADER", "true") == "true",

		HCaptchaSecret:  env.OrDefault("HCAPTCHA_SECRET", ""),
		HCaptchaSitekey: env.OrDefault("HCAPTCHA_SITEKEY", ""),

		PaymentsAPIAddress: env.OrDefault("PAYMETSAPI_ADDRESS", "localhost"),
	}
}
