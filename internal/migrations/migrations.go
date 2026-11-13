package migrations

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var FS embed.FS

func Start(ctx context.Context, pool *pgxpool.Pool) error {
	sql, err := FS.ReadFile("0_init.sql")
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, string(sql), pgx.QueryExecModeSimpleProtocol)
	return err
}
