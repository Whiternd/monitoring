package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "monitoring/forwarder"
	"monitoring/handler"
	"monitoring/storage"
)

func main() {
	cfg := handler.Config{
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
		ForwardURLs: strings.Split(os.Getenv("FORWARD_URLS"), ","),
		ListenAddr:  ":8080",
	}

	var db *storage.Storage
	var err error
	for i := 0; i < 10; i++ {
		db, err = storage.NewPostgres(cfg.PostgresDSN)
		if err == nil {
			break
		}
		log.Printf("Postgres not ready, retrying in 2s: %v", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect to Postgres after retries: %v", err)
	}

	err = db.InitSchema()
	if err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	app := handler.NewApp(db, cfg.ForwardURLs)

	http.HandleFunc("/ingest", app.HandleIngest)
	log.Printf("Listening on %s", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, nil))
}
