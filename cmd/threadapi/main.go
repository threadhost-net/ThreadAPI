package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/threadhost-net/ThreadAPI/internal/config"
	"github.com/threadhost-net/ThreadAPI/internal/httpapi"
	"github.com/threadhost-net/ThreadAPI/internal/migrations"
	"github.com/threadhost-net/ThreadAPI/platform"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	time.Local = time.UTC

	cfg := config.New()

	stack := newStack(cfg)
	defer stack.Redis.Close()
	defer stack.DB.Close()

	migrationCtx, migrationCancel := context.WithTimeout(ctx, 10*time.Second)
	if err := migrations.Start(migrationCtx, stack.DB); err != nil {
		migrationCancel()
		log.Println("Не удалось создать структуру БД: ", err)
		return
	}
	migrationCancel()

	handler := httpapi.InitRouter(stack, cfg.Domain)
	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println("Сервер выключен принудительно: ", err)
		}
		log.Println("Сервер выключен")
	}()

	ln, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Println("Не удалось занять адрес: ", err)
		return
	}
	log.Println("Сервер запущен: ", ln.Addr())

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.Serve(ln)
	}()

	select {
	case <-ctx.Done():
		log.Println("Сервер выключается...")
	case err := <-srvErr:
		log.Println("Сервер упал: ", err)
		return
	}
}

func newStack(cfg *config.Config) *platform.Stack {
	return &platform.Stack{
		Redis:  initRedis(cfg.RedisHost, cfg.RedisPassword),
		Mailer: newMailer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpFrom),
		DB:     initDBPool(cfg.DatabaseURL),
		Domain: cfg.Domain,
	}
}

func newMailer(host string, port int, username, password, from string) *platform.Mailer {
	return &platform.Mailer{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}
}

func initRedis(redisHost, redisPassword string) *platform.MultiRedis {
	managerRedis := platform.NewMultiRedis().WithDefault(redisHost, redisPassword)
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
	managerDBPool, err := platform.NewDBPool(databaseURL)
	if err != nil {
		log.Fatalln("postgres init failed:", err)
	}
	return managerDBPool
}
