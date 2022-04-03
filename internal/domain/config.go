package domain

import validation "github.com/go-ozzo/ozzo-validation"

// AppConfig stores all configuration of the application.
type AppConfig struct {
	CowboyName      string `env:"COWBOY_NAME"`
	CowboyTableName string `env:"COWBOY_TABLE_NAME"`
}

// Validate - validates configuration of the application.
func (cfg AppConfig) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.CowboyName, validation.Required),
		validation.Field(&cfg.CowboyTableName, validation.Required),
	)
}

// CronConfig stores all configuration of the cron application.
type CronConfig struct {
	CowboyTableName string `env:"COWBOY_TABLE_NAME"`
}

// Validate - validates configuration of the application.
func (crn CronConfig) Validate() error {
	return validation.ValidateStruct(&crn,
		validation.Field(&crn.CowboyTableName, validation.Required),
	)
}
