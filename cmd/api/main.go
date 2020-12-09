package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"test-document/pkg/config"
	"test-document/pkg/docs"
	"test-document/pkg/storage"
)

func main() {
	cfg := config.New()

	db, err := sqlx.Connect("postgres", cfg.GetString("DB"))
	if err != nil {
		os.Exit(1)
	}

	docsRepo := storage.NewRepository(db)
	docsSvc := docs.NewService(docsRepo)

	r := initRouter(docsSvc)
	srv := initServer(cfg, r)

	go func() {
		fmt.Println("server is running")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Errorf("error listen and server: %w", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Errorf("error server shutdown: %w", err)
	}

	os.Exit(0)
}

func initRouter(docsSvc docs.Service) *mux.Router {
	r := mux.NewRouter()

	versionRout := r.PathPrefix("/api").Subrouter()

	versionRout.HandleFunc("/hello", MakeHandler()).Methods(http.MethodGet)

	r.PathPrefix("/document").Handler(docs.NewNewsHandler(versionRout, docsSvc))

	return r
}

func initServer(cfg *viper.Viper, r *mux.Router) *http.Server {
	return &http.Server{
		Addr:    cfg.GetString("LISTEN"),
		Handler: r,
	}
}

func MakeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]bool{"hello": true})
		if err != nil {
			log.Println(err)
		}
	}
}
