package usecases

import (
	"strings"

	"github.com/chanchai9666/aider"
	"github.com/gofiber/fiber/v3"
	"github.com/jinzhu/copier"

	"template-fiber-v3/internal/entities/models"
	"template-fiber-v3/internal/entities/schemas"
)

func (s *userRequest) FindUsers(c fiber.Ctx, req *schemas.FindUsersReq) ([]models.Users, error) {
	data, err := s.repo.FindUsers(req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// FindUsersByEmail ค้นหาผู้ใช้งานตามอีเมล
func (s *userRequest) FindUsersByEmail(c fiber.Ctx, req *schemas.FindUsersByEmailReq) (*models.Users, error) {
	if req.Email == "" {
		return nil, aider.NewError(aider.ErrBadRequest, "กรุณาระบุอีเมล")
	}

	data, err := s.FindUsers(c, &schemas.FindUsersReq{Email: req.Email})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, aider.NewError(aider.ErrNotFound, "ไม่พบข้อมูล")
	}
	return &data[0], nil
}

// FindUsersAll ค้นหาผู้ใช้งานทั้งหมด
func (s *userRequest) FindUsersAll() ([]models.Users, error) {
	data, err := s.repo.FindUsers(&schemas.FindUsersReq{})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateUsers สร้างผู้ใช้งาน
func (s *userRequest) CreateUsers(c fiber.Ctx, req *schemas.AddUsers) error {
	return s.repo.CreateUsers(req)
}

// UpdateUsers อัพเดตผู้ใช้งาน
func (s *userRequest) UpdateUsers(c fiber.Ctx, req *schemas.AddUsers) error {
	return s.repo.UpdateUser(req)
}

// DeleteUsers ลบผู้ใช้งาน
func (s *userRequest) DeleteUsers(c fiber.Ctx, req *schemas.AddUsers) error {
	return s.repo.DeletedUser(req.UserId)
}

// Login ล็อกอิน
func (s *userRequest) Login(c fiber.Ctx, req *schemas.LoginReq) (*schemas.LoginResp, error) {
	// ดึงข้อมูลผู้ใช้ตาม UserID
	result, err := s.FindUsersByEmail(c, &schemas.FindUsersByEmailReq{Email: req.Email})
	if err != nil {
		return nil, err
	}

	// ตรวจสอบการเข้าสู่ระบบ
	if result.UserId == 0 || result.Password == "" {
		return nil, aider.NewError(aider.ErrInternal, "ข้อมูลไม่ครบถ้วน")
	}

	// ตรวจสอบรหัสผ่าน
	if !aider.CheckPassword(req.Password, result.Password) {
		return nil, aider.NewError(aider.ErrInternal, "รหัสผ่านไม่ถูกต้อง")
	}

	// คัดลอกข้อมูลผู้ใช้
	var userLogin schemas.LoginResp
	var userData schemas.UserResp
	if err := copier.Copy(&userData, result); err != nil {
		return nil, err
	}

	levelVal := []string{}
	for _, v := range result.Level {
		levelVal = append(levelVal, v.Level)
	}
	strLvl := strings.Join(levelVal, ",")
	userLogin.User = userData
	userLogin.User.Level = levelVal

	// สร้าง JWT
	token, err := s.repo.NewJwt(&schemas.JwtReq{
		UserId:  aider.ToString(result.UserId),
		Name:    result.Name,
		SurName: result.SurName,
		Email:   result.Email,
		Level:   strLvl,
	})
	if err != nil {
		return nil, err
	}

	userLogin.AccessToken = "Bearer " + token
	aider.DDD(userLogin)
	return &userLogin, nil
}

// RefreshToken รีเฟรช JWT
func (s *userRequest) RefreshToken(c fiber.Ctx, req *schemas.RefreshTokenReq) (*schemas.LoginResp, error) {
	// ดึงข้อมูลผู้ใช้ตาม UserID
	result, err := s.FindUsersByEmail(c, &schemas.FindUsersByEmailReq{Email: req.Email})
	if err != nil {
		return nil, err
	}
	// ตรวจสอบความถูกต้อง
	if req.UserId != result.UserId {
		return nil, aider.NewError(aider.ErrInternal, "ข้อมูลไม่ครบถ้วน")
	}

	// คัดลอกข้อมูลผู้ใช้
	var userLogin schemas.LoginResp
	var userData schemas.UserResp
	if err := copier.Copy(&userData, result); err != nil {
		return nil, err
	}

	levelVal := []string{}
	for _, v := range result.Level {
		levelVal = append(levelVal, v.Level)
	}
	strLvl := strings.Join(levelVal, ",")
	userLogin.User = userData
	userLogin.User.Level = levelVal

	// สร้าง JWT
	token, err := s.repo.NewJwt(&schemas.JwtReq{
		UserId:  aider.ToString(result.UserId),
		Name:    result.Name,
		SurName: result.SurName,
		Email:   result.Email,
		Level:   strLvl,
	})
	if err != nil {
		return nil, err
	}

	userLogin.AccessToken = "Bearer " + token
	aider.DDD(userLogin)
	return &userLogin, nil

}

func (s *userRequest) FindByUserID(c fiber.Ctx, req *schemas.FindByUserIDReq) (*models.Users, error) {
	if req.UserID == "" {
		return nil, aider.NewError(aider.ErrBadRequest, "กรุณาระบุ UserID")
	}

	data, err := s.FindUsers(c, &schemas.FindUsersReq{UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, aider.NewError(aider.ErrNotFound, "ไม่พบข้อมูล")
	}
	return &data[0], nil
}
