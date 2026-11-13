package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/threadhost-net/ThreadAPI/internal/applications"
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

	stack := platform.Stack{
		Redis:          platform.InitRedis(cfg.RedisHost, cfg.RedisPassword),
		Mailer:         platform.NewMailer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpFrom),
		DB:             platform.InitDBPool(cfg.DatabaseURL),
		Domain:         cfg.Domain,
		HCaptchaSecret: cfg.HCaptchaSecret,
	}
	defer stack.Redis.Close()
	defer stack.DB.Close()

	migrationCtx, migrationCancel := context.WithTimeout(ctx, 10*time.Second)
	defer migrationCancel()
	if err := migrations.Start(migrationCtx, stack.DB); err != nil {
		log.Fatalln("Не удалось создать структуру БД: ", err)
	}

	handler := httpapi.InitRouter(applications.NewToolBox(&stack), cfg.Domain)
	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Сервер выключен принудительно: ", err)
	}
	log.Println("Сервер выключен")
}
