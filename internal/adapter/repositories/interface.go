package repositories

import (
	"gorm.io/gorm"

	"template-fiber-v3/configs"
	"template-fiber-v3/internal/entities/models"
	"template-fiber-v3/internal/entities/schemas"
)

// func NewConstRepository(db *gorm.DB, config *configs.Config, userId string) ConstRepository {
// 	return &constDB{
// 		constRequest: &constRequest{db: db},
// 		userId:       userId, // กำหนดค่าให้กับ Name
// 		// config:       config,
// 	}
// }

func NewUsersRepository(db *gorm.DB, config *configs.Config) UsersRepository {
	return &dbPool{
		database: database{db: db}, // ใช้ชื่อฟิลด์เพื่อกำหนดค่า
		config:   config,
	}
}

// func NewProductRepository(db *gorm.DB) ProductRepository {
// 	return &constDB{
// 		constRequest: &constRequest{db: db},
// 	}
// }

type UsersRepository interface {
	CreateUsers(req *schemas.AddUsers) error                     //เพิ่มข้อมูล Users
	FindUsers(req *schemas.FindUsersReq) ([]models.Users, error) //ค้นหาข้อมูล Users
	UpdateUser(req *schemas.AddUsers) error                      //อัพเดตข้อมูล Users
	DeletedUser(userID uint64) error                             //ลบข้อมูล Users
	NewJwt(req *schemas.JwtReq) (string, error)                  //สร้าง Jwt Token
}

// type ConstRepository interface {
// 	Create(req *schemas.ConfigConstant) error                                         //เพิ่ม ค่าคงที่
// 	Update(req *schemas.ConfigConstant) error                                         //แก้ไข ค่าคงที่
// 	Delete(id, group string) error                                                    //ลบค่าคงที่
// 	FindPage(req *schemas.ConfigConstant) (*Pagination[models.ConfigConstant], error) //ค้าหาค่าคงที่แบบแบ่งหน้า
// 	FindAll(req *schemas.ConfigConstant) ([]models.ConfigConstant, error)             //ค้นหาค่าคงที่ทั้งหมด
// }

// type ProductRepository interface {
// 	AddProduct(req *schemas.ProductReq) error
// }
