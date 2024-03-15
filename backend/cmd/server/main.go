package main

import (
	"backend/cmd/server/config"
	_ "backend/contracts"
	"backend/internal/infra/firefly"
	"backend/internal/infra/mysql"
	"backend/internal/middleware"
	"backend/internal/service/item"
	http2 "backend/internal/transport/http"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

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
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatalf("Failed to load timezone for DB: %s", err.Error())
		return exitError
	}
	cc := mysql2.Config{
		DBName:    "kaleido",
		User:      "yurie",
		Passwd:    cfg.MysqlPassword,
		Addr:      "docker.for.mac.localhost:3306",
		Net:       "tcp",
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}
	db, err := sql.Open("mysql", cc.FormatDSN())
	if err != nil {
		log.Fatalf("Failed to prepare DB: %s", err.Error())
		return exitError
	}
	defer db.Close()

	dbClient := mysql.New(db)
	httpClient := http.DefaultClient
	httpUrl1, err := url.Parse(cfg.FireflyBaseUrlUserOne)
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", cfg.FireflyBaseUrlUserOne, err.Error())
		return exitError
	}
	httpUrl2, err := url.Parse(cfg.FireflyBaseUrlUserTwo)
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", cfg.FireflyBaseUrlUserTwo, err.Error())
		return exitError
	}
	httpUrl3, err := url.Parse(cfg.FireflyBaseUrlUserThree)
	if err != nil {
		log.Fatalf("Failed to parse url (%s): %s", cfg.FireflyBaseUrlUserThree, err.Error())
		return exitError
	}
	fireflyClient := firefly.New(httpUrl1, httpUrl2, httpUrl3, httpClient)
	itemService := item.New(fireflyClient, dbClient)
	httpServer := http2.New(itemService)

	log.Println("Creating NFT pool...")
	if err := fireflyClient.CreatePool(context.Background()); err != nil {
		log.Fatalf("Failed to create pool for NFT via Firefly: %s", err.Error())
		return exitError
	}

	log.Println("Setting up HTTP server...")
	r := mux.NewRouter()
	r.Use(middleware.SetUserID)
	r.HandleFunc("/items/list", httpServer.ListItem).Methods("POST")
	r.HandleFunc("/items/buy", httpServer.PurchaseItem).Methods("POST")
	r.HandleFunc("/items/get", httpServer.GetItem).Methods("GET")

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
