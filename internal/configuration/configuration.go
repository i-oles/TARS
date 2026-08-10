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

type Sender struct {
	Host      string
	Port      int
	Login     string
	Password  string
	Signature string
}

type DoctorReminder struct {
	RecipientEmail string
	Interval       Duration
}

type Tasks struct {
	DoctorReminder DoctorReminder
}

type Configuration struct {
	ListenAddress       string
	DBPath              string
	ReadTimeout         Duration
	WriteTimeout        Duration
	ContextTimeout      Duration
	AuthSecret          string
	LogBusinessErrors   bool
	LogConfig           bool
	Sender              Sender
	DomainAddr          string
	BaseRequestTmplPath string
	MockEmailSender     bool
	Tasks               Tasks
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

	if cfg.Sender.Login == "" || cfg.Sender.Password == "" {
		return nil,
			errors.New("provide envs for notifier")
	}

	return cfg, nil
}

func loadEnvs(cfg *Configuration) {
	if login := os.Getenv("SENDER_LOGIN"); login != "" {
		cfg.Sender.Login = login
	}

	if password := os.Getenv("SENDER_PASSWORD"); password != "" {
		cfg.Sender.Password = password
	}

	if authSecret := os.Getenv("AUTH_SECRET"); authSecret != "" {
		cfg.AuthSecret = authSecret
	}

	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}

	if drReminderRecipientEmail := os.Getenv("DR_REMINDER_RECIPIENT_EMAIL"); drReminderRecipientEmail != "" {
		cfg.Tasks.DoctorReminder.RecipientEmail = drReminderRecipientEmail
	}
}
