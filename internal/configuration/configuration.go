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
	Host      string
	Port      int
	Login     string
	Password  string
	Signature string
}

type Tasks struct {
	DoctorReminder DoctorReminder
	CeneoCatcher   CeneoCatcher
}

type DoctorReminder struct {
	RecipientEmail string
	RefID          string
	Interval       Duration
}

type CeneoCatcher struct {
	MaxLimitPrice int
	ProductID     int
	Interval      Duration
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
	Mailer            Mailer
	Scheduler         Scheduler
	DomainAddr        string
	BaseEmailTmplPath string
	MockMailer        bool
	Tasks             Tasks
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

	if cfg.Mailer.Login == "" || cfg.Mailer.Password == "" {
		return nil,
			errors.New("provide envs for sender")
	}

	return cfg, nil
}

func loadEnvs(cfg *Configuration) {
	if login := os.Getenv("SENDER_LOGIN"); login != "" {
		cfg.Mailer.Login = login
	}

	if password := os.Getenv("SENDER_PASSWORD"); password != "" {
		cfg.Mailer.Password = password
	}

	if authSecret := os.Getenv("AUTH_SECRET"); authSecret != "" {
		cfg.AuthSecret = authSecret
	}

	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.DBPath = dbPath
	}

	if drReminderRecipientEmail := os.Getenv(
		"DR_REMINDER_RECIPIENT_EMAIL",
	); drReminderRecipientEmail != "" {
		cfg.Tasks.DoctorReminder.RecipientEmail = drReminderRecipientEmail
	}

	if refID := os.Getenv("DR_REMINDER_REF_ID"); refID != "" {
		cfg.Tasks.DoctorReminder.RefID = refID
	}
}
