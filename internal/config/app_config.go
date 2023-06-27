package config

import "github.com/spf13/viper"

type AppConfig struct {
	DbMigrationPath    string `mapstructure:"DB_MIGRATION_PATH"`
	DbMigrationVersion int    `mapstructure:"DB_MIGRATION_VERSION"`

	DbURL             string `mapstructure:"DB_URL"`
	DbUseSSL          bool   `mapstructure:"DB_USE_SSL"`
	DbMaxConn         int    `mapstructure:"DB_MAX_CONN"`
	DbAcquireTimeout  int    `mapstructure:"DB_ACQUIRE_TIMEOUT"`
	DbValidationQuery string `mapstructure:"DB_VALIDATION_QUERY"`
}

func LoadAppConfig() (AppConfig, error) {

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)

	viper.BindEnv("DB_MIGRATION_PATH")
	viper.BindEnv("DB_MIGRATION_VERSION")
	viper.BindEnv("DB_URL")
	viper.BindEnv("DB_USE_SSL")
	viper.BindEnv("DB_MAX_CONN")
	viper.BindEnv("DB_ACQUIRE_TIMEOUT")
	viper.BindEnv("DB_VALIDATION_QUERY")

	viper.SetDefault("DB_VALIDATION_QUERY", "SELECT 1")
	viper.SetDefault("DB_MAX_CONN", 10)
	viper.SetDefault("DB_USE_SSL", true)
	viper.SetDefault("DB_ACQUIRE_TIMEOUT", 5)

	var appConfig AppConfig
	err := viper.Unmarshal(&appConfig)
	if err != nil {
		return appConfig, err
	}

	return appConfig, nil
}
