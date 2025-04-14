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

	// s.App.Use(requestid.New())
	// s.App.Use(middleware.LogRequestMiddleware()) // Custom middleware for logging requests
	s.App.Use(logger.New(logger.Config{
		Format: logger.ConfigDefault.Format,
	}))
	jwtAuth := middleware.AuthMiddleware(s.config.JwtSECRETKEY) // JWT Authentication middleware

	// Group base path
	api := s.App.Group("/api")

	// Swagger UI (optional)
	if s.config.SwagStatus {
		api.Get("/docs/*", swagger.HandlerDefault)
	}

	// Init dependencies
	userRepo := repositories.NewUsersRepository(s.db.DB(), s.config)
	userService := usecases.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService) // <- เปลี่ยนชื่อให้สอดคล้อง

	// Health check & login
	api.Get("/health", s.healthHandler)
	api.Post("/login", userHandler.Login)

	// Auth group (require JWT)
	auth := api.Group("/auth", jwtAuth)
	auth.Post("/refreshToken", userHandler.RefreshToken)

	// Users group
	users := api.Group("/users", jwtAuth)
	users.Post("/", userHandler.CreateUsers)                 // Create user
	users.Get("/", userHandler.FindUser)                     // Find users by filters
	users.Get("/all", userHandler.FindAllUsers)              // Find all
	users.Get("/:user_id<\\d+>", userHandler.FindByUserID)   // Find by ID (Match แค่ตัวเลขเท่านั้น เช่น /api/users/123)
	users.Delete("/:user_id<\\d+>", userHandler.DeleteUsers) // Delete user

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
