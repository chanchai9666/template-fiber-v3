package handlers

import (
	"github.com/gofiber/fiber/v3"

	"template-fiber-v3/internal/adapter"
	_ "template-fiber-v3/internal/entities/models"
	"template-fiber-v3/internal/entities/schemas"
	"template-fiber-v3/internal/usecases"
)

type userEndPoint struct {
	service usecases.UsersService
}

func NewUserEndPoint(service usecases.UsersService) *userEndPoint {
	return &userEndPoint{service: service}
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
func (en *userEndPoint) FindUser(c fiber.Ctx) error {
	return adapter.RespJson(c, en.service.FindUsers, &schemas.FindUsersReq{})
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
func (en *userEndPoint) FindUsersByUserId(c fiber.Ctx) error {
	return adapter.RespJson(c, en.service.FindUsersByEmail, &schemas.FindUsersByEmailReq{})
}

// @Tags Users
// @Summary ค้นหา User ทั้งหมด
// @Description Show User ทั้งหมด
// @Accept  json
// @Produce  json
// @Param Accept-Language header string false "(en, th)" default(th)
// @Success 200 {object} []models.Users
// @Failure 400 {object} schemas.HTTPError
// @Router /api/users/usersAll [get]
// @Security ApiKeyAuth
func (en *userEndPoint) FindUserAll(c fiber.Ctx) error {
	return adapter.RespJsonNoReq(c, en.service.FindUsersAll)
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
// @Router /api/users/createUsers [post]
// @Security ApiKeyAuth
func (en *userEndPoint) CreateUsers(c fiber.Ctx) error {
	return adapter.RespSuccess(c, en.service.CreateUsers, &schemas.AddUsers{})
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
func (en *userEndPoint) UpdateUsers(c fiber.Ctx) error {
	return adapter.RespSuccess(c, en.service.UpdateUsers, &schemas.AddUsers{})
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
// @Router /api/users/deleteUsers/{user_id} [delete]
// @Security ApiKeyAuth
func (en *userEndPoint) DeleteUsers(c fiber.Ctx) error {
	return adapter.RespSuccess(c, en.service.DeleteUsers, &schemas.AddUsers{})
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
func (en *userEndPoint) Login(c fiber.Ctx) error {
	return adapter.RespJson(c, en.service.Login, &schemas.LoginReq{})
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
func (en *userEndPoint) RefreshToken(c fiber.Ctx) error {
	return adapter.RespJson(c, en.service.RefreshToken, &schemas.RefreshTokenReq{})
}
