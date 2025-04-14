package repositories

import (
	"errors"
	"fmt"
	"reflect"
	"template-fiber-v3/configs"

	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

type mainDB struct {
	*database
	userId string
	// config *configs.Config
}

type dbPool struct {
	database
	config *configs.Config
}

type constDB struct {
	*database
	userId string
	// config *configs.Config
}

// กำหนดให้ใหม่ให้ repo ของ user ใหม่เนื่องจากมีการรับค่า string ที่จำเป็นต้องใช้งาน
type userDB struct {
	database
	config *configs.Config
}

// Insert function for inserting data into the database
//
//	func Insert[T any](db *gorm.DB, data T) error {
//		return db.Create(&data).Error // ใช้ &data เพื่อส่ง pointer
//	}

// Insert คือฟังก์ชัน generic สำหรับ insert ข้อมูลเข้า database โดยรองรับกรณีต่อไปนี้:
//
// ✅ &Struct      → pointer ไปยัง struct → GORM จะอัปเดต field อัตโนมัติ เช่น ID
// ✅ Struct       → struct ปกติ          → แปลงเป็น pointer ภายในฟังก์ชัน
// ✅ &[]Struct    → pointer ไปยัง slice  → ใช้ insert หลายรายการ + ได้รับค่า ID คืน
// ✅ &[]*Struct   → pointer ไปยัง slice ของ pointer → รองรับได้เช่นกัน
//
// ❌ []Struct     → ไม่รองรับ เพราะ GORM ไม่สามารถอัปเดตค่า ID กลับเข้า slice เดิมได้
// ❌ ประเภทอื่น ๆ เช่น int, string, map → ไม่รองรับ
//
// ✅ คำแนะนำ: ถ้าต้องการให้ GORM อัปเดตค่าอัตโนมัติ เช่น ID ที่เพิ่มใหม่ → ควรส่ง pointer
func Insert[T any](db *gorm.DB, data T) error {
	v := reflect.ValueOf(data)
	k := v.Kind()

	switch k {
	case reflect.Ptr:
		// ถ้าเป็น pointer → อาจเป็น *Struct หรือ *[]Struct
		elemKind := v.Elem().Kind()
		if elemKind == reflect.Struct || elemKind == reflect.Slice {
			// ✅ กรณีที่รองรับ: ใช้งาน data ตรง ๆ ได้เลย
			// - *Struct   → insert รายการเดียว, คืนค่า ID ได้
			// - *[]Struct หรือ *[]*Struct → insert หลายรายการ, คืนค่า ID ได้
			return db.Create(data).Error
		}
	case reflect.Struct:
		// ✅ ถ้าเป็น struct เดี่ยว (ไม่ใช่ pointer) → แปลงเป็น pointer ก่อน
		// → เพื่อให้ GORM อัปเดต ID และ field อัตโนมัติ
		ptr := v.Addr().Interface()
		return db.Create(ptr).Error
	case reflect.Slice:
		// ❌ ถ้าเป็น slice ธรรมดา → ไม่รองรับ
		// เพราะ Go ส่งแบบ pass by value → ค่า ID จะไม่ถูกอัปเดตกลับ
		return errors.New("กรุณาส่ง pointer ไปยัง slice เช่น &[]User{} เพื่อให้สามารถอัปเดตค่าอัตโนมัติ (เช่น ID) ได้")
	}

	// ❌ กรณีอื่น ๆ เช่น map, int, string → ไม่รองรับ
	return errors.New("ไม่รองรับประเภทข้อมูลนี้ใน Insert")
}

func Updates[T any](db *gorm.DB, data *T) error {
	tx := db.Model(data).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func Delete[T any](db *gorm.DB, data *T) error {

	tx := db.Select("is_active", "deleted_user").Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	tx = db.Delete(data)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("delete error : record not found")
	}
	return nil
}

func UpdateInterface[T any](db *gorm.DB, model *T, data map[string]interface{}) error {
	return db.Model(model).Updates(data).Error
}

func Find[T any](db *gorm.DB, data *T) error {
	if err := db.Find(data).Error; err != nil {
		return err
	}
	return nil
}

func Count[T any](db *gorm.DB, data *T) int64 {
	var count int64
	if err := db.Model(data).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// สำหรับ Transaction
func Transaction(db *gorm.DB, handler func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	fmt.Println("======")
	fmt.Println("Transaction Begin")
	fmt.Println("======")

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("======")
			fmt.Println("Transaction Panic -> Rollback")
			fmt.Println("======")
			panic(r) // ส่ง panic ต่อไป
		} else if err != nil {
			tx.Rollback()
			fmt.Println("======")
			fmt.Println("Transaction Error -> Rollback")
			fmt.Println("======")
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = commitErr
				tx.Rollback()
				fmt.Println("======")
				fmt.Println("Transaction Commit Error -> Rollback")
				fmt.Println("======")
			} else {
				fmt.Println("======")
				fmt.Println("Transaction Commit")
				fmt.Println("======")
			}
		}
	}()

	// main handler
	err = handler(tx)
	return
}

// Create create record database
func Create(db *gorm.DB, i interface{}) error {
	if err := db.Create(i).Error; err != nil {
		return err
	}
	return nil
}
