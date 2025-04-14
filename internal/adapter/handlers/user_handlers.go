package handlers

import (
	"github.com/gofiber/fiber/v3"

	"template-fiber-v3/internal/adapter"
	_ "template-fiber-v3/internal/entities/models"
	"template-fiber-v3/internal/entities/schemas"
	"template-fiber-v3/internal/usecases"
)

type userHandler struct {
	service usecases.UsersService
}

func NewUserHandler(service usecases.UsersService) *userHandler {
	return &userHandler{service: service}
}

// @Tags Users
// @Summary ค้นหา User ตามเงื่อนไข
// @Description Show User ตามเงื่อนไข
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request query schemas.FindUsersReq false " request body "
// @Success 200 {object} []models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users [get]
// @Security ApiKeyAuth
func (h *userHandler) FindUser(c fiber.Ctx) error {
	return adapter.RespJson(c, h.service.FindUsers, &schemas.FindUsersReq{})
}

// @Tags Users
// @Summary ค้นหา User ตาม UserId
// @Description Show User ตาม UserId
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param user_id path string true "User ID"
// @Success 200 {object} models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/{user_id} [get]
// @Security ApiKeyAuth
func (h *userHandler) FindByUserID(c fiber.Ctx) error {
	return adapter.RespJson(c, h.service.FindByUserID, &schemas.UserIDReq{})
}

// @Tags Users
// @Summary ค้นหา User ทั้งหมด
// @Description Show User ทั้งหมด
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} []models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/all [get]
// @Security ApiKeyAuth
func (h *userHandler) FindAllUsers(c fiber.Ctx) error {
	return adapter.RespJsonNoReq(c, h.service.FindUsersAll)
}

// @Tags Users
// @Summary เพิ่มข้อมูล User
// @Description เพิ่มข้อมูล User
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body schemas.AddUsers false " request body "
// @Success 200 {object} schemas.HTTPError
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users [post]
// @Security ApiKeyAuth
func (h *userHandler) CreateUsers(c fiber.Ctx) error {
	return adapter.RespSuccess(c, h.service.CreateUsers, &schemas.AddUsers{})
}

// @Tags Users
// @Summary แก้ไขข้อมูล User
// @Description แก้ไขข้อมูล User
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body schemas.AddUsers false " request body "
// @Success 200 {object} schemas.HTTPError
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/updateUsers [post]
// @Security BearerAuth
func (h *userHandler) UpdateUsers(c fiber.Ctx) error {
	return adapter.RespSuccess(c, h.service.UpdateUsers, &schemas.AddUsers{})
}

// @Tags Users
// @Summary ลบข้อมูล User
// @Description ลบข้อมูล User
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param user_id path string true "User ID"
// @Success 200 {object} schemas.HTTPError
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/{user_id} [delete]
// @Security ApiKeyAuth
func (h *userHandler) DeleteUsers(c fiber.Ctx) error {
	return adapter.RespSuccess(c, h.service.DeleteUsers, &schemas.UserIDReq{})
}

// @Tags Auth
// @Summary Login เข้าใช้งาน
// @Description Login เข้าใช้งานสำหรับขอ token
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body schemas.LoginReq false " request body "
// @Success 200 {object} schemas.LoginResp
// @Failure 400 {object} schemas.HTTPError
// @Router /api/login [post]
func (h *userHandler) Login(c fiber.Ctx) error {
	return adapter.RespJson(c, h.service.Login, &schemas.LoginReq{})
}

// @Tags Auth
// @Summary ขอ Token เข้าใช้งานระบบใหม่
// @Description Refresh เพื่อขอ Token เข้าใช้งานระบบใหม่
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Accept-Language header string false "(en, th)" default(th)
// @Param request body schemas.RefreshTokenReq false "request body"
// @Success 200 {object} schemas.LoginResp
// @Failure 400 {object} schemas.HTTPError
// @Router /api/auth/refreshToken [post]
func (h *userHandler) RefreshToken(c fiber.Ctx) error {
	return adapter.RespJson(c, h.service.RefreshToken, &schemas.RefreshTokenReq{})
}
