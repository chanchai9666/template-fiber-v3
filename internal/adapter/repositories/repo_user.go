package repositories

import (
	"fmt"

	"github.com/chanchai9666/aider"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"template-fiber-v3/internal/entities/models"
	"template-fiber-v3/internal/entities/schemas"
	"template-fiber-v3/internal/pkg/safety"
)

func (r *dbPool) CreateUsers(req *schemas.AddUsers) error {
	var user models.Users
	if err := copier.Copy(&user, req); err != nil {
		return fmt.Errorf("failed to copy user data: %w", err)
	}

	return Transaction(r.db, func(tx *gorm.DB) error {
		return Insert(r.db, &user)
	})
}
func (r *dbPool) FindUsers(req *schemas.FindUsersReq) ([]models.Users, error) {

	var allUsers []models.Users
	tx := r.db

	pagination := &Pagination[models.Users]{
		Sort: "email asc",
	}
	err := tx.Preload("Level").Scopes(
		WhereIsActive(),
		WhereName(req.Name),
		WhereSurName(req.SurName),
		WhereEmail(req.Email),
		WhereUserId(req.UserId),
		Paginate(r.db, models.Users{}, pagination),
	).Find(&allUsers).Error
	if err != nil {
		return nil, err
	}
	pagination.Rows = allUsers
	return allUsers, nil
}
func (r *dbPool) UpdateUser(req *schemas.AddUsers) error {
	var users models.Users
	if err := copier.Copy(&users, req); err != nil {
		return fmt.Errorf("failed to copy user data: %w", err)
	}

	return Transaction(r.db, func(tx *gorm.DB) error {
		que := r.db.Select("name", "sur_name").Scopes(WhereUserId(aider.ToString(req.UserId)))
		return Updates(que, &users)
	})
}
func (r *dbPool) DeletedUser(userID uint64) error {
	return Transaction(r.db, func(d *gorm.DB) error {
		var UserUpdate models.Users
		// สร้างตัวแปรสำหรับอัปเดต
		active := 0
		UserUpdate.IsActive = &active
		UserUpdate.UserId = userID
		// ลบผู้ใช้
		if err := Delete(
			r.db.Scopes(WhereUserId(aider.ToString(userID))), &UserUpdate); err != nil {
			return err
		}
		return nil
	})
}

func (r *dbPool) NewJwt(req *schemas.JwtReq) (string, error) {
	jwt, err := safety.GenerateJWT(r.config.JwtSECRETKEY, &safety.JwtConst{
		UserId:   req.UserId,
		Name:     req.Name,
		SurName:  req.SurName,
		Email:    req.Email,
		Level:    req.Level,
		SafetyId: "xxxxxxx",
	})
	if err != nil {
		return "nil", err
	}
	return jwt, nil
}
