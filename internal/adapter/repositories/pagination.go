package repositories

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty" query:"limit"`
	Page       int    `json:"page,omitempty" query:"page"`
	Sort       string `json:"sort,omitempty" query:"sort"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       []T    `json:"rows" swaggertype:"array,object"`
}

func (p *Pagination[T]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination[T]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 30
	}
	return p.Limit
}

func (p *Pagination[T]) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[T]) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id desc"
	}
	return p.Sort
}

// ส่วนการแบ่งหน้า
func Paginate[T any](db *gorm.DB, models T, pagination *Pagination[T]) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(models).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}

	//Ex ตัวอย่างการใช้งาน
	/*
		func GetUsers(db *gorm.DB) (*Pagination[User], error) {
			var users []User
			pagination := &Pagination[User]{
					Limit: 10,
					Page:  1,
			}

			err := db.Scopes(Paginate(User{}, pagination, db)).Find(&users).Error
			if err != nil {
					return nil, err
			}

			pagination.Rows = users
			return pagination, nil
		}
	*/

}

// แบ่งหน้าแบบ RawSql
func PaginateRawSQL(db *gorm.DB, rawSQL string, pagination *Pagination[any]) (*gorm.DB, error) {
	// คำนวณ total rows
	var totalRows int64
	err := db.Raw("SELECT COUNT(*) FROM (" + rawSQL + ") as count_query").Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	// สร้าง paginated query
	paginatedSQL := rawSQL + fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.GetLimit(), pagination.GetOffset())

	// เพิ่ม ORDER BY ถ้าจำเป็น
	if pagination.GetSort() != "" {
		paginatedSQL += " ORDER BY " + pagination.GetSort()
	}

	return db.Raw(paginatedSQL), nil

	//Ex ตัวอย่างการใช้งาน
	/*
		func GetPaginatedUsersWithPosts(db *gorm.DB, pagination *Pagination[UserWithPosts]) ([]UserWithPosts, error) {
			rawSQL := `
					SELECT u.id, u.name, u.email, p.title as post_title
					FROM users u
					LEFT JOIN posts p ON u.id = p.user_id
					WHERE u.active = true
			`

			query, err := PaginateRawSQL(db, pagination, rawSQL)
			if err != nil {
					return nil, err
			}

			var usersWithPosts []UserWithPosts
			err = query.Scan(&usersWithPosts).Error
			if err != nil {
					return nil, err
			}

			pagination.Rows = usersWithPosts
			return usersWithPosts, nil
		}
	*/

}
