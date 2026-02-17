package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/antnzr/gatemulator/config"
	app "github.com/antnzr/gatemulator/internal/app/gatemulator"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()
			isostring := "2006-01-02T15:04:05.999Z07:00"
			a.Value = slog.StringValue(t.Format(isostring))
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	conf, err := config.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}

	db, err = sql.Open("sqlite3", filepath.Join(dir, "gatemulator.db"))
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	defer db.Close()

	migarations := `
		CREATE TABLE IF NOT EXISTS subscriptions (
			id UUID PRIMARY KEY,
			subscription_id UUID NOT NULL,
			subscriber_id UUID,
			transport VARCHAR(20),
			phone VARCHAR(20),
			state VARCHAR(50),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(subscription_id, phone)
		);

		CREATE TABLE IF NOT EXISTS subscribers (
			id UUID PRIMARY KEY,
			title TEXT NOT NULL UNIQUE,
			token TEXT NOT NULL UNIQUE,
			webhook_url TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS message_files (
			id UUID PRIMARY KEY,
			sha1 TEXT NOT NULL,
			mime_type TEXT NOT NULL,
			message_id UUID NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS message_file_sha1 ON message_files (sha1);
	`
	_, err = db.Exec(migarations)
	if err != nil {
		slog.Error("Failed to create table", slog.Any("err", err))
	}

	srv := app.NewServer(conf, db)

	ec := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		ec <- srv.Run(context.Background())
	}()

	// Waits for an internal error that shutdowns the server.
	// Otherwise, wait for a SIGINT or SIGTERM and tries to shutdown the server gracefully.
	// After a shutdown signal, HTTP requests taking longer than the specified grace period are forcibly closed.
	select {
	case err = <-ec:
	case <-ctx.Done():
		haltCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		srv.Shutdown(haltCtx)
		stop()
		err = <-ec
	}

	slog.Info("Buy 🖖🏾 %v", slog.Any("err", err))
}
