package handler

import (
	"io"
	"log"
	"net/http"

	"monitoring/forwarder"
	"monitoring/storage"
)

type Config struct {
	PostgresDSN string
	ForwardURLs []string
	ListenAddr  string
}

type App struct {
	DB        *storage.Storage
	Forwarder *forwarder.Forwarder
}

func NewApp(db *storage.Storage, urls []string) *App {
	return &App{
		DB:        db,
		Forwarder: forwarder.NewForwarder(urls),
	}
}

func (a *App) HandleIngest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received: %s", body)

	err = a.DB.Insert(body)
	if err != nil {
		log.Printf("Failed to insert into Postgress: %v", err)
	} else {
		log.Printf("Successfully inserted into Postgress")
	}

	go a.Forwarder.Forward(body)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
