package platform

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Stack struct {
	Redis          *MultiRedis
	Mailer         *Mailer
	DB             *pgxpool.Pool
	Domain         string
	HCaptchaSecret string
}

func InitRedis(redisHost, redisPassword string) *MultiRedis {
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

func InitDBPool(databaseURL string) *pgxpool.Pool {
	managerDBPool, err := NewDBPool(databaseURL)
	if err != nil {
		log.Fatalln("postgres init failed:", err)
	}
	return managerDBPool
}
