package main

import (
	"backend/cmd/server/config"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello World!")
}

const (
	exitOK = iota
	exitError
	shutdownWait = time.Second * 30
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to setup config: %s", err.Error())
	}

	os.Exit(Run(cfg))
}

func Run(cfg *config.Config) int {
	log.Println("Setting up DB...")
	db, err := sql.Open("mysql", cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("Failed to prepare DB: %s", err.Error())
		return exitError
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to health-check DB: %s", err.Error())
		return exitError
	}

	log.Println("Setting up HTTP server...")
	r := mux.NewRouter()
	r.HandleFunc("/", TestHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP server: %s", err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), shutdownWait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Gracefully shutting down...")
	return exitOK
}
