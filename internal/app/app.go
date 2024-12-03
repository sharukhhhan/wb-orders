package app

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"wb-l-zero/config"
	v1 "wb-l-zero/internal/controller/http/v1"
	"wb-l-zero/internal/kafka"
	"wb-l-zero/internal/repository"
	"wb-l-zero/internal/service"
	"wb-l-zero/pkg/logger"
	"wb-l-zero/pkg/validator"
)

func Run(configPath string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurations set up
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// logger
	logger.SetupLogrus(cfg.Log.Level)
	kafkaLogger, err := logger.NewFileLogger("./logs/kafka.log")

	// validator
	valid := validator.NewCustomValidator()

	// Connecting to Postgres
	log.Info("Connecting postgres...")
	pg, err := pgx.Connect(ctx, cfg.PG.URL)
	if err != nil {
		log.Fatal(fmt.Errorf("error connecting postgres: %w", err))
	}
	defer pg.Close(ctx)

	// Running Migrations
	log.Info("Running migrations...")
	err = RunMigrations(cfg.PG.URL, cfg.PG.MigrationPath)
	if err != nil {
		log.Debug(fmt.Errorf("error running migrations: %w", err))
	}

	//  Cache
	bigCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(cfg.Cache.TTL))
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repository.NewRepository(pg, bigCache)

	// Service
	log.Info("Initializing services")
	services := service.NewService(repositories)

	// Handler
	log.Info("Initializing handlers and routes...")
	handler := echo.New()
	handler.Validator = valid
	v1.NewRouter(handler, services)

	// Kafka Consumer
	log.Info("Starting Kafka consumer...")
	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topic, services, valid, kafkaLogger)
	if err != nil {
		log.Fatalf("Error initializing Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := kafkaConsumer.Start(ctx)
		if err != nil {
			kafkaLogger.Errorf("Kafka error: %v", err)
		} else {
			log.Info("Kafka Consumer stopped gracefully.")
		}
	}()

	// HTTP server
	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: handler,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server stopped with error: %v", err)
		} else {
			log.Info("HTTP server stopped gracefully.")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	<-stop
	log.Info("Shutting down...")

	// Shutdown Kafka Consumer
	cancel()

	// Shutdown HTTP Server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Error during server shutdown: %v", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Info("All components stopped gracefully.")
}
