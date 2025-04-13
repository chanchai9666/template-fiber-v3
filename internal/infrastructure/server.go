package infrastructure

import (
	"context"
	"template-fiber-v3/configs"
	"template-fiber-v3/internal/infrastructure/database"

	"github.com/gofiber/fiber/v3"
)

type FiberServer struct {
	App    *fiber.App
	db     database.Service
	config *configs.Config
}

func New(configs *configs.Config) *FiberServer {
	// Create database configuration

	dbConfig := database.Config{
		Host:     configs.DBHost,
		Port:     configs.DBPort,
		Database: configs.DBDatabase,
		Username: configs.DBUsername,
		Password: configs.DBPassword,
		Schema:   configs.DBSchema,
	}

	// Create database service
	db := database.New("default", dbConfig)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ArcZed",
			AppName:      "ArcZed",
		}),
		db:     db,
		config: configs,
	}

	return server
}

// Listen starts the server on the specified address
func (s *FiberServer) Listen(addr string) error {
	return s.App.Listen(addr, fiber.ListenConfig{
		EnablePrefork: true,
	})
}

// ShutdownWithContext gracefully shuts down the server
func (s *FiberServer) ShutdownWithContext(ctx context.Context) error {
	return s.App.ShutdownWithContext(ctx)
}
