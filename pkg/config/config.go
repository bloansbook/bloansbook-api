package config

import (
	"log"

	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Supabase SupabaseConfig
	Resend   ResendConfig
	IDGen    IDGenConfig
}

type AppConfig struct {
	Env            string
	Port           string
	AllowedOrigins string
}

type DatabaseConfig struct {
	URL string
}

type SupabaseConfig struct {
	URL       string
	AnonKey   string
	JWTSecret string
}

type ResendConfig struct {
	APIKey string
}

type IDGenConfig struct {
	StaffPrefix    string
	CustomerPrefix string
	SupplierPrefix string
}

var ApplicationConfig *Config

func Load() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("%s", sysmsg.NoEnvFile)
	}

	ApplicationConfig = &Config{
		App: AppConfig{
			Env:            viper.GetString("APP_ENV"),
			Port:           viper.GetString("APP_PORT"),
			AllowedOrigins: viper.GetString("ALLOWED_ORIGINS"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
		Supabase: SupabaseConfig{
			URL:       viper.GetString("SUPABASE_URL"),
			AnonKey:   viper.GetString("SUPABASE_ANON_KEY"),
			JWTSecret: viper.GetString("JWT_SECRET"),
		},
		Resend: ResendConfig{
			APIKey: viper.GetString("RESEND_API_KEY"),
		},
		IDGen: IDGenConfig{
			StaffPrefix:    viper.GetString("STAFF_ID_PREFIX"),
			CustomerPrefix: viper.GetString("CUSTOMER_ID_PREFIX"),
			SupplierPrefix: viper.GetString("SUPPLIER_ID_PREFIX"),
		},
	}
}
