package models

import "github.com/shopspring/decimal"

type Product struct {
	AuditLog
	ProductId   int             `json:"product_id" gorm:"type:int;primaryKey;autoIncrement"` //id สินค้า
	ProductName string          `json:"product_name" gorm:"type:varchar(255)"`               //ชื่อสินค้า
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(10,2)"`                     //ราคาสินค้า
	UnitId      string          `json:"unit_id" gorm:"type:varchar(50)"`                     //หน่วยสินค้า
	CategoryId  int             `json:"category_id" gorm:"type:int"`                         //หมวดหมู่สินค้า
	Description string          `json:"description" gorm:"type:text"`                        //รายละเอียดสินค้า
	ImageURL    string          `json:"image_url" gorm:"type:varchar(255)"`                  //URL รูปสินค้า
}
