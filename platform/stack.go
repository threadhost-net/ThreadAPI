package platform

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/threadhost-net/ThreadAPI/internal/config"
)

type Stack struct {
	Redis           *MultiRedis
	Mailer          *Mailer
	DB              *pgxpool.Pool
	Domain          string
	HCaptchaSecret  string
	HCaptchaSitekey string
}

func NewStack(cfg *config.Config) *Stack {
	return &Stack{
		Redis:           initRedis(cfg.RedisHost, cfg.RedisPassword),
		Mailer:          newMailer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpFrom),
		DB:              initDBPool(cfg.DatabaseURL),
		Domain:          cfg.Domain,
		HCaptchaSecret:  cfg.HCaptchaSecret,
		HCaptchaSitekey: cfg.HCaptchaSitekey,
	}
}

func newMailer(host string, port int, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func initRedis(redisHost, redisPassword string) *MultiRedis {
	managerRedis := NewMultiRedis().WithDefault(redisHost, redisPassword)
	managerRedis.Add("sessions", 0)
	managerRedis.Add("email_verify", 1)
	managerRedis.Add("password_reset", 2)

	if errs := managerRedis.Errors(); len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
		log.Fatalln("redis init failed")
	}

	return managerRedis
}

func initDBPool(databaseURL string) *pgxpool.Pool {
	managerDBPool, err := NewDBPool(databaseURL)
	if err != nil {
		log.Fatalln("postgres init failed:", err)
	}
	return managerDBPool
}
