package infrastructure

import (
	"template-fiber-v3/docs"
	"template-fiber-v3/internal/adapter/handlers"
	"template-fiber-v3/internal/adapter/repositories"
	"template-fiber-v3/internal/pkg/middleware"
	"template-fiber-v3/internal/usecases"
	"time"

	swagger "github.com/chanchai9666/swagger-fiber-v3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

// RegisterFiberRoutes registers all routes for the Fiber server
func (s *FiberServer) RegisterFiberRoutes() {
	// กำหนดรายละเอียดของส่วน auth Bearer
	// @securityDefinitions.apikey ApiKeyAuth
	// @name Authorization
	// @in ใส่ค่า Bearer เว้นวรรคและตามด้วย TOKEN  ex(Bearer ?????????)
	// END กำหนดค่าให้ swagger
	// =======================================================
	// เพิ่ม middleware สำหรับการเข้าถึง Swagger UI ด้วยควบคุมสิทธิ์

	docs.SwaggerInfo.Host = s.config.SwagHost
	docs.SwaggerInfo.Title = s.config.SwagTitle
	docs.SwaggerInfo.Description = s.config.SwagDescription
	docs.SwaggerInfo.Version = s.config.SwagVersion
	docs.SwaggerInfo.BasePath = s.config.SwagBasePath
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Enable CORS
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowCredentials: false,
		MaxAge:           int(time.Hour * 4), // 4 ชั่วโมง
	}))
	s.App.Use(recover.New())

	s.App.Use(requestid.New())
	s.App.Use(middleware.LogRequestMiddleware()) // Custom middleware for logging requests
	s.App.Use(logger.New(logger.Config{
		Format: logger.ConfigDefault.Format,
	}))
	jwtAuth := middleware.AuthMiddleware(s.config.JwtSECRETKEY) // JWT Authentication middleware

	api := s.App.Group("/api")
	if s.config.SwagStatus {
		api.Get("/docs/*", swagger.HandlerDefault) // Swagger UI
	}

	auth := api.Group("/auth", jwtAuth)
	_ = auth
	userRep := repositories.NewUsersRepository(s.db.DB(), s.config) // User repository
	userService := usecases.NewUserService(userRep)                 // User service
	userEndpoint := handlers.NewUserEndPoint(userService)           // User endpoint

	api.Get("/", userEndpoint.FindUser)
	api.Get("/health", s.healthHandler)                   // Health check database
	api.Post("/login", userEndpoint.Login)                // Login endpoint
	auth.Post("/refreshToken", userEndpoint.RefreshToken) // Refresh token endpoint

	api.Post("/users", userEndpoint.CreateUsers)               // Create user
	api.Get("/users", userEndpoint.FindUser)                   // Find user by query parameters
	api.Get("/users/:user_id", userEndpoint.FindUsersByUserId) // Find user by user ID
	api.Get("/users/usersAll", userEndpoint.FindUserAll)       // Find all users

}

// healthHandler handles the health check endpoint
func (s *FiberServer) healthHandler(c fiber.Ctx) error {
	// Get database health status
	dbHealth := s.db.Health()
	// Return health status
	return c.JSON(fiber.Map{
		"database": dbHealth,
	})
}
