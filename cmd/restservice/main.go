package main

import (
	"fmt"
	"github.com/MarouaneMan/github-api/internal/config"
	"github.com/MarouaneMan/github-api/internal/restservice"
	"github.com/MarouaneMan/github-api/kvstore"
	"github.com/MarouaneMan/github-api/middleware"
	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/kelseyhightower/envconfig"
	"net/http"
	"os"
	"time"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")

	// Init config
	cfg := &config.Config{}
	{
		err := envconfig.Process("", cfg)
		if err != nil {
			log.WithError(err).Error("Fail to initialize configuration")
			os.Exit(1)
		}
	}

	// Instantiate a new key value store
	store := kvstore.NewInMemoryStore(30*time.Second, 30*time.Minute)

	// Spawn fetcher job to periodically pull Github data
	// !! Usually this block needs to run in its own process/dedicated node when using redis as a backend store
	//{
	//	ctx := logger.ToCtx(context.Background(), log)
	//	cronScheduler := gocron.NewScheduler(time.UTC)
	//	cronScheduler.Every(cfg.FetchIntervalHours).Hours().StartImmediately().Do(func() {
	//		fetcher.Run(ctx, cfg, store)
	//	})
	//	cronScheduler.StartAsync()
	//}

	log.Info("Initializing routes")
	router := handlers.NewRouter(log)
	router.Use(middleware.NewResponseCachingMiddleware(store, store))
	router.HandleFunc("/ping", restservice.PongHandler).Methods("GET", "POST")
	router.HandleFunc("/repos", restservice.ReposHandler(store)).Methods("GET")
	router.HandleFunc("/stats", restservice.StatsHandler(store)).Methods("GET")

	log = log.WithField("port", cfg.Port)
	log.Info("Listening...")
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", cfg.Port), router) // TODO: listen on 0.0.0.0
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}
