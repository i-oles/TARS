package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tkanos/gonfig"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("could not unmarshal duration string: %w", err)
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("could not parse duration: %w", err)
	}

	d.Duration = parsed

	return nil
}

type Mailer struct {
	Login     string
	Password  string
	Signature string
}

type Mailers struct {
	Host  string
	Port  int
	Owner Mailer
	Tars  Mailer
}

type Scheduler struct {
	Interval Duration
}

type Configuration struct {
	ListenAddress     string
	DBPath            string
	ReadTimeout       Duration
	WriteTimeout      Duration
	ContextTimeout    Duration
	AuthSecret        string
	LogBusinessErrors bool
	LogConfig         bool
	Mailers           Mailers
	Scheduler         Scheduler
	DomainAddr        string
	BaseEmailTmplPath string
	MockMailer        bool
}

func (c *Configuration) Pretty() string {
	cfgPretty, _ := json.MarshalIndent(c, "", "  ")

	return string(cfgPretty)
}

func GetConfig(cfgPath string) (*Configuration, error) {
	cfg := &Configuration{}

	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found, using environment variables...")
	}

	configFileName := os.Getenv("CONFIG")

	cfgFinalPath := filepath.Join(cfgPath, configFileName+".json")

	err = gonfig.GetConf(cfgFinalPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not load configuration: %s", err.Error())
	}

	loadEnvs(cfg)

	if cfg.Mailers.Owner.Login == "" || cfg.Mailers.Owner.Password == "" {
		return nil,
			errors.New("provide envs for owner mailer")
	}

	if cfg.Mailers.Tars.Login == "" || cfg.Mailers.Tars.Password == "" {
		return nil,
			errors.New("provide envs for tars mailer")
	}

	return cfg, nil
}

func loadEnvs(cfg *Configuration) {
	if login := os.Getenv("OWNER_MAILER_LOGIN"); login != "" {
		cfg.Mailers.Owner.Login = login
	}

	if password := os.Getenv("OWNER_MAILER_PASSWORD"); password != "" {
		cfg.Mailers.Owner.Password = password
	}

	if login := os.Getenv("TARS_MAILER_LOGIN"); login != "" {
		cfg.Mailers.Tars.Login = login
	}

	if password := os.Getenv("TARS_MAILER_PASSWORD"); password != "" {
		cfg.Mailers.Tars.Password = password
	}

	if authSecret := os.Getenv("AUTH_SECRET"); authSecret != "" {
		cfg.AuthSecret = authSecret
	}

	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}
}
