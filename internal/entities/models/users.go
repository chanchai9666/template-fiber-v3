package models

// ข้อมูล User
type Users struct {
	AuditLog
	UserId      uint64        `json:"user_id" gorm:"primaryKey;autoIncrement"`                      //ไอดี ของผู้ใช้งาน
	Email       string        `json:"email" gorm:"type:varchar(100);unique;comment:อีเมล"`          //อีเมล
	Password    string        `json:"password" gorm:"type:varchar(255);comment:รหัสผ่าน"`           //ชื่อ โปรไฟล์
	Name        string        `json:"name" gorm:"type:varchar(50);comment:ชื่อ"`                    //ชื่อ
	SurName     string        `json:"sur_name" gorm:"type:varchar(50);comment:นามสกุล"`             //นามสกุล
	PhoneNumber *string       `json:"phone_number" gorm:"type:varchar(20);comment:หมายเลขโทรศัพท์"` //หมายเลขโทรศัพท์
	IdCard      *string       `json:"id_card" gorm:"type:varchar(30);comment:รหัสบัตรประจำตัว"`     //รหัสบัตรประจำตัว
	Level       []UsersLevels `gorm:"foreignKey:UserID;references:UserId"`                          //ความสัมพันธ์กับ UserLevels
}

type UsersLevels struct {
	AuditLog
	ID     uint   `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID uint64 `gorm:"column:user_id;index;type:varchar(50);not null"` // ต้องกำหนดประเภทและขนาดให้ตรงกับ Users.UserId
	Level  string `gorm:"column:level;type:varchar(50);not null"`
}
