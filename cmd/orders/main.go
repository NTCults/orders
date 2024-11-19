package main

import (
	"context"
	"net/http"
	"orders/internal/config"
	"orders/internal/query"
	"orders/internal/service"
	"orders/internal/store"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	if cfg.Debug {
		log.Info("Starting app in debug mode")
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.JSONFormatter{})

	store := store.NewStore(cfg)
	defer store.Close()

	ordersQuery, err := query.NewNatsOrdersQuery(store, cfg)
	if err != nil {
		log.WithField(config.NATS_URL, cfg.NatsURL).Fatal("Unable to connect to nats")
	}

	service := service.NewService(store)

	srv := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      service.Router,
	}

	go func() {
		log.WithField("Addr", srv.Addr).Info("Starting server.")
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	<-c

	ordersQuery.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info("Shutting down")

	os.Exit(0)
}
