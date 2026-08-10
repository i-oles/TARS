package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"main/internal/application/task"
	"main/internal/configuration"
	"main/internal/domain/repositories"
	"main/internal/domain/requester"
	sqliteRepo "main/internal/infrastructure/repository/sqlite"
	"main/internal/infrastructure/requester/gmail"
	"main/internal/infrastructure/scheduler"
	"main/internal/infrastructure/sender"
	gmailSender "main/internal/infrastructure/sender/gmail"
	"main/internal/infrastructure/sender/memory"
	"main/internal/interfaces/http/html/handlers/emails"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Components struct {
	tasksRepo repositories.ITasks
	requester requester.IRequester
	database  *gorm.DB
	storage   *memory.Storage
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("failed to loading configuration", slog.String("err", err.Error()))
		os.Exit(1)
	}

	components, err := buildComponents(cfg)
	if err != nil {
		slog.Error("failed to build components", slog.String("err", err.Error()))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctxTimeout time.Duration) {
		schedulerCtx, cancel := context.WithTimeout(ctx, ctxTimeout)
		defer cancel()

		doctorReminderTask := task.NewDoctorReminder(
			components.requester,
			cfg.Tasks.DoctorReminder.RecipientEmail,
		)

		scheduler := scheduler.New(cfg.Tasks.DoctorReminder.Interval.Duration, doctorReminderTask)

		err := scheduler.Run(schedulerCtx)
		if err != nil {
			slog.Error("failed to run scheduler", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}(cfg.ContextTimeout.Duration)

	router := setupRouter(
		components.storage,
		cfg,
	)

	srv := &http.Server{
		Addr:              cfg.ListenAddress,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadTimeout.Duration,
		ReadTimeout:       cfg.ReadTimeout.Duration,
		WriteTimeout:      cfg.WriteTimeout.Duration,
	}

	runServer(srv, cfg)
}

func loadConfig() (*configuration.Configuration, error) {
	cfg, err := configuration.GetConfig("./config")
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	if cfg.LogConfig {
		slog.Info(cfg.Pretty())
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	return cfg, nil
}

func buildComponents(cfg *configuration.Configuration) (Components, error) {
	database, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return Components{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Successfully connected to database")

	err = database.AutoMigrate(
		&sqliteRepo.SQLTask{},
	)
	if err != nil {
		return Components{}, fmt.Errorf("failed to migrate database: %w", err)
	}

	var sender sender.IEmailSender

	sender = gmailSender.New(
		cfg.Sender.Host,
		cfg.Sender.Port,
		cfg.Sender.Login,
		cfg.Sender.Password,
	)

	memoryStorage := memory.Storage{
		Views: make([]string, 0),
	}

	if cfg.MockEmailSender {
		sender = memory.NewSender(&memoryStorage)
	}

	tasksRepo := sqliteRepo.NewTasksRepo(database)

	requester := gmail.NewRequester(
		sender,
		cfg.Sender.Login,
		cfg.Sender.Signature,
		cfg.BaseRequestTmplPath,
	)

	return Components{
		tasksRepo: tasksRepo,
		requester: requester,
		database:  database,
		storage:   &memoryStorage,
	}, nil
}

func setupRouter(
	storage *memory.Storage,
	cfg *configuration.Configuration,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/")

	emailsHander := emails.NewHandler(storage)

	{
		// testing
		if cfg.MockEmailSender {
			api.GET("/emails", emailsHander.Handle)
		}
	}

	// API
	// authMiddleware := middleware.Auth(cfg.AuthSecret)
	//
	// createTaskHandler := createtask.NewHandler(tasksRepo)
	// listTasksHandler := listbookings.NewHandler(tasksRepo)
	//
	// {
	// 	api.POST("/api/v1/tasks", authMiddleware, createTaskHandler.Handle)
	// 	api.GET("/api/v1/tasks", authMiddleware, listTasksHandler.Handle)
	// }

	return router
}

func runServer(srv *http.Server, cfg *configuration.Configuration) {
	go func() {
		slog.Info("Starting server...", slog.String("address", cfg.ListenAddress))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ContextTimeout.Duration)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("err", err.Error()))
	}

	slog.Info("Server stopped")
}
