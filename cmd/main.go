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

	"main/internal/application/email"
	"main/internal/application/task/ceneocatcher"
	"main/internal/application/task/doctorreminder"
	"main/internal/configuration"
	"main/internal/domain/contracts"
	"main/internal/infrastructure/mailer/gmail"
	"main/internal/infrastructure/mailer/memory"
	sqliteRepo "main/internal/infrastructure/repository/sqlite"
	"main/internal/infrastructure/scheduler"
	"main/internal/interfaces/http/api/errs"
	"main/internal/interfaces/http/api/handlers/updatetask"
	"main/internal/interfaces/http/html/handlers/emails"
	"main/internal/interfaces/middleware"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Components struct {
	tasksRepo           contracts.ITasks
	emailComposer       email.Composer
	mailer              email.IMailer
	errHandler          errs.IErrorHandler
	memoryMailerStorage *memory.Storage
	database            *gorm.DB
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

	go func() {
		doctorReminderTask := doctorreminder.New(
			components.emailComposer,
			components.mailer,
			components.tasksRepo,
			cfg.Tasks.DoctorReminder.RecipientEmail,
			cfg.Tasks.DoctorReminder.RefID,
			cfg.Tasks.DoctorReminder.Interval.Duration,
		)

		ceneoCatcherTask := ceneocatcher.New(
			components.emailComposer,
			components.mailer,
			components.tasksRepo,
			215,
			198572377,
			1*time.Minute,
		)

		scheduler := scheduler.New(
			cfg.Scheduler.Interval.Duration,
			doctorReminderTask,
			ceneoCatcherTask,
		)

		if err := scheduler.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("failed to run scheduler", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	router := setupRouter(
		components.memoryMailerStorage,
		components.errHandler,
		components.tasksRepo,
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

	var mailer email.IMailer

	mailer = gmail.NewMailer(
		cfg.Mailer.Host,
		cfg.Mailer.Port,
		cfg.Mailer.Login,
		cfg.Mailer.Password,
	)

	memoryMailerStorage := memory.Storage{
		Views: make([]string, 0),
	}

	if cfg.MockMailer {
		mailer = memory.NewMailer(&memoryMailerStorage)
	}

	tasksRepo := sqliteRepo.NewTasksRepo(database)

	emailComposer := email.NewComposer(
		mailer,
		cfg.Mailer.Login,
		cfg.Mailer.Signature,
		cfg.BaseEmailTmplPath,
	)

	errHandler := errs.NewErrorHandler()

	return Components{
		tasksRepo:           tasksRepo,
		emailComposer:       *emailComposer,
		mailer:              mailer,
		database:            database,
		memoryMailerStorage: &memoryMailerStorage,
		errHandler:          errHandler,
	}, nil
}

func setupRouter(
	mailerMemoryStorage *memory.Storage,
	errHandler errs.IErrorHandler,
	tasksRepo contracts.ITasks,
	cfg *configuration.Configuration,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/")

	{
		// testing
		if cfg.MockMailer {
			emailsHander := emails.NewHandler(mailerMemoryStorage)

			api.GET("/emails", emailsHander.Handle)
		}
	}

	authMiddleware := middleware.Auth(cfg.AuthSecret)

	updateTaskHandler := updatetask.NewHandler(tasksRepo, errHandler)
	api.PATCH("/api/v1/tasks/:task_id", authMiddleware, updateTaskHandler.Handle)

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
