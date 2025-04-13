package repositories

import (
	"fmt"

	"gorm.io/gorm"
)

// ฟังก์ชันตรวจสอบและแปลงค่าจาก data[0]
func getFirstValue(data ...any) (string, bool) {
	if len(data) == 0 {
		return "", false // ไม่มีข้อมูล
	}
	value := fmt.Sprintf("%v", data[0])
	if value == "" {
		return "", false // ค่าที่แปลงเป็น string ว่าง
	}
	return value, true
}

// where amount>1000 [แบบฟิกค่า]
func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	return db.Where("amount > ?", 1000)
}

// where pay_mode=card [แบบฟิกค่า]
func PaidWithCreditCard(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode = ?", "card")
}

// where pay_mode=cod [แบบฟิกค่า]
func PaidWithCod(db *gorm.DB) *gorm.DB {
	return db.Where("pay_mode = ?", "cod")
}

func WhereTable(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 1️⃣ ตรวจสอบจำนวนพารามิเตอร์ ป้องกัน panic
		if len(data) < 3 {
			return db
		}

		// 2️⃣ ตรวจสอบประเภทของ tableName และ fieldName
		tableName, ok1 := data[0].(string)
		fieldName, ok2 := data[1].(string)
		if !ok1 || !ok2 {
			return db // คืนค่า db ถ้าไม่ใช่ string
		}

		// 3️⃣ ตรวจสอบประเภทของ value
		value := data[2]

		// ✅ ถ้า value เป็น nil ให้ใช้ IS NULL แทน
		if value == nil {
			condition := fmt.Sprintf("%s.%s IS NULL", tableName, fieldName)
			return db.Where(condition)
		}

		// ✅ รองรับประเภทข้อมูลที่ใช้กับ GORM ได้
		switch value.(type) {
		case string, int, float64, float32, uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64, bool:
			// ✅ ใช้ได้
		default:
			return db // ❌ ป้องกันค่าที่ไม่รองรับ
		}

		// 4️⃣ ใช้ tableName และ fieldName อย่างปลอดภัย
		condition := fmt.Sprintf("%s.%s = ?", tableName, fieldName)
		return db.Where(condition, value)
	}
}

// 0=table 1=field 2=value
func WhereLikeTable(data ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(data) < 3 {
			return db
		}
		tableName := data[0]
		fieldName := data[1]
		value := data[2]
		condition := fmt.Sprintf("%s.%s LIKE ?", tableName, fieldName)
		return db.Where(condition, "%"+value+"%")
	}
}

func WhereIsActive(status ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "is_active"
		isActive := "1" // Default to active

		// Assign default active status if no argument is provided
		if len(status) > 0 && status[0] != "" {
			isActive = status[0]
		}

		// Check if isActive is set to inactive ("0")
		if isActive == "0" {
			db = db.Unscoped() // Ignore soft-deleted records
		} else if isActive == "-1" {
			return db.Unscoped()
		}

		// Additional conditions based on multiple arguments
		if len(status) > 1 {
			return WhereTable(isActive, fieldName, status[1])(db)
		}

		return db.Where(fieldName+" = ?", isActive)
	}
}

// where user_id
func WhereUserId(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "user_id"
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → เรียกใช้ WhereTable
		if len(data) > 1 {
			return WhereTable(value, fieldName, data[1])(db)
		}
		return db.Where(fieldName+" = ?", value)
	}
}

// where const_id
func WhereConstId(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "const_id"
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → เรียกใช้ WhereTable
		if len(data) > 1 {
			return WhereTable(value, fieldName, data[1])(db)
		}
		return db.Where(fieldName+" = ?", value)
	}
}

// where group_id
func WhereGroupId(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "group_id"
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → เรียกใช้ WhereTable
		if len(data) > 1 {
			return WhereTable(value, fieldName, data[1])(db)
		}
		return db.Where(fieldName+" = ?", value)
	}
}

// where email
func WhereEmail(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "email"
		// ใช้ฟังก์ชัน getFirstValue แทนการเขียนโค้ดซ้ำ
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → เรียกใช้ WhereTable
		if len(data) > 1 {
			return WhereTable(value, fieldName, data[1])(db)
		}
		return db.Where(fieldName+" = ?", value)
	}
}

func WhereEmail2(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "email"
		// ใช้ฟังก์ชัน getFirstValue แทนการเขียนโค้ดซ้ำ
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → เรียกใช้ WhereTable
		if len(data) > 1 {
			return WhereTable(value, fieldName, data[1])(db)
		}
		// ค้นหาตามปกติ
		return db.Where(fieldName+" = ?", value)
	}
}

// where name
func WhereName(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "name"
		// ใช้ฟังก์ชัน getFirstValue เพื่อรับค่า
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → แปลง data[1] ให้เป็น string และเรียกใช้ WhereLikeTable
		if len(data) > 1 {
			fieldValue, ok := data[1].(string) // แปลง data[1] ให้เป็น string
			if !ok {
				return db // ถ้าไม่สามารถแปลงเป็น string ได้จะคืนค่า db
			}
			return WhereLikeTable(value, fieldName, fieldValue)(db)
		}
		// ค้นหาตามปกติ
		condition := fieldName + " LIKE ?"
		return db.Where(condition, "%"+value+"%")
	}
}

// where sur_name
func WhereSurName(data ...any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fieldName := "sur_name"
		// ใช้ฟังก์ชัน getFirstValue เพื่อรับค่า
		value, valid := getFirstValue(data...)
		if !valid {
			return db
		}
		// กรณีมีพารามิเตอร์มากกว่า 1 ตัว → แปลง data[1] ให้เป็น string และเรียกใช้ WhereLikeTable
		if len(data) > 1 {
			fieldValue, ok := data[1].(string) // แปลง data[1] ให้เป็น string
			if !ok {
				return db // ถ้าไม่สามารถแปลงเป็น string ได้จะคืนค่า db
			}
			return WhereLikeTable(value, fieldName, fieldValue)(db)
		}
		// ค้นหาตามปกติ
		condition := fieldName + " LIKE ?"
		return db.Where(condition, "%"+value+"%")
	}
}
