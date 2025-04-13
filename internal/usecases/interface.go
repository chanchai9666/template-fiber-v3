package usecases

import (
	"template-fiber-v3/internal/adapter/repositories"
	"template-fiber-v3/internal/entities/models"
	"template-fiber-v3/internal/entities/schemas"

	"github.com/gofiber/fiber/v3"
)

type UsersService interface {
	FindUsers(c fiber.Ctx, req *schemas.FindUsersReq) ([]models.Users, error)              // ค้นหาผู้ใช้งาน
	FindUsersAll() ([]models.Users, error)                                                 // ค้นหาผู้ใช้งานทั้งหมด
	CreateUsers(c fiber.Ctx, req *schemas.AddUsers) error                                  // สร้างผู้ใช้งาน
	FindUsersByEmail(c fiber.Ctx, req *schemas.FindUsersByEmailReq) (*models.Users, error) // ค้นหาผู้ใช้งานตามอีเมล
	UpdateUsers(c fiber.Ctx, req *schemas.AddUsers) error                                  // อัพเดตผู้ใช้งาน
	DeleteUsers(c fiber.Ctx, req *schemas.AddUsers) error                                  // ลบผู้ใช้งาน
	Login(c fiber.Ctx, req *schemas.LoginReq) (*schemas.LoginResp, error)                  // ล็อกอิน
	RefreshToken(c fiber.Ctx, req *schemas.RefreshTokenReq) (*schemas.LoginResp, error)    // รีเฟรช JWT
}

// struct (โครงสร้างข้อมูล) ที่ใช้เป็นตัว “implement” interface UsersService
type userRequest struct {
	repo repositories.UsersRepository
}

// เป็น constructor function สำหรับ UsersService
func NewUserService(repo repositories.UsersRepository) UsersService {
	return &userRequest{
		repo: repo,
	}
}
