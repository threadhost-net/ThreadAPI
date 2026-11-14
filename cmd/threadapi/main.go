package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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

	stack := platform.NewStack(cfg)
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
