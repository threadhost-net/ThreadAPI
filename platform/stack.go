package platform

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Stack struct {
	Redis  *MultiRedis
	Mailer *Mailer
	DB     *pgxpool.Pool
	Domain string
}
