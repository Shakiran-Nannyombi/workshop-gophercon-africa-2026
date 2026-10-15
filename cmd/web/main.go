// Command web serves Karibu Guide.
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"workshop/internal/handlers"
)

func main() {
	port := flag.Int("port", envInt("PORT", 8080), "port to listen on")
	dev := flag.Bool("dev", os.Getenv("KARIBU_DEV") == "1", "serve /static and /data from disk")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := run(logger, *port, *dev); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func run(logger *slog.Logger, port int, dev bool) error {
	rate, _ := strconv.ParseFloat(os.Getenv("KARIBU_RATE_KES"), 64)
	h, err := handlers.New(logger, handlers.Options{
		Dev:     dev,
		BaseURL: os.Getenv("KARIBU_BASE_URL"),
		RateKES: rate,
	})
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           handlers.Logged(logger, h),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	logger.Info("Karibu Guide listening", "url", "http://localhost:"+strconv.Itoa(port), "dev", dev)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func envInt(key string, fallback int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return v
	}
	return fallback
}
