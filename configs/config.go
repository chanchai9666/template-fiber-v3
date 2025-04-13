package configs

import (
	"log"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port   int    `env:"PORT" envDefault:"8080"`
	AppEnv string `env:"APP_ENV" envDefault:"local"`

	// Database
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBDatabase string `env:"DB_DATABASE" envDefault:"test"`
	DBUsername string `env:"DB_USERNAME" envDefault:"admin"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"1234"`
	DBSchema   string `env:"DB_SCHEMA" envDefault:"public"`

	//JWT
	JwtSECRETKEY string `env:"SECRET_KEY" envDefault:"the_secret_key"`

	// Swagger Info
	SwagStatus      bool   `env:"SWAG_INFO_STATUS" envDefault:"false"`
	SwagHost        string `env:"SWAG_INFO_HOST" envDefault:"localhost:8080"`
	SwagTitle       string `env:"SWAG_INFO_TITLE" envDefault:"MyApp"`
	SwagDescription string `env:"SWAG_INFO_DESCRIPTION" envDefault:"API Description"`
	SwagVersion     string `env:"SWAG_INFO_VERSION" envDefault:"1.0"`
	SwagContactName string `env:"SWAG_INFO_CONTACT_NAME" envDefault:"My Company"`
	SwagBasePath    string `env:"SWAG_INFO_BASE_PATH" envDefault:"/"`
}

func LoadConfig() (*Config, error) {
	var config Config

	envPath := filepath.Join("configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	// ใช้ env.Parse เพื่อโหลดค่า environment variables ลงใน struct
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
