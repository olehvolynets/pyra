package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/joho/godotenv/autoload"

	"pyra/pkg/db"
	"pyra/pkg/log"
	"pyra/pkg/server"
	"pyra/pkg/session"

	"pyra/internal/api"
)

func main() {
	rootLogger := log.NewLogger()

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		done()
		if r := recover(); r != nil {
			rootLogger.Error("application panic", "panic", r)
			os.Exit(1)
		}
	}()

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
	if err != nil {
		rootLogger.Error("failed to parse port value", "error", err)
		os.Exit(1)
	}
	srv, err := server.New(server.WithPort(uint(port)), server.WithLogger(rootLogger))
	if err != nil {
		rootLogger.Error("failed to create an http server", "error", err)
		os.Exit(1)
	}

	dbConf := db.NewConfig("pgx")
	dbPool, err := db.New(ctx, dbConf, rootLogger)
	if err != nil {
		rootLogger.Error("failed to create a DB pool", "error", err)
		os.Exit(1)
	}

	mux := api.Mux(dbPool, rootLogger)

	staticFs := http.FileServerFS(os.DirFS("./public/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", staticFs))

	var h http.Handler = mux
	h = server.PanicRecovery(h)
	session.SetupSessionStore()
	h = server.Session(h)
	h = server.Logger(rootLogger, h)

	if err = srv.Start(ctx, h); err != nil {
		rootLogger.Error("failed to start a server", "error", err)
	}
}
