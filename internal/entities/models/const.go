package models

type ConfigConstant struct {
	AuditLog
	ConstID   string `json:"const_id" gorm:"type:varchar(50);primaryKey;comment:ไอดีค่าคงที่;"`
	GroupID   string `json:"group_id" gorm:"type:varchar(50);primaryKey;index;comment:ไอดีกลุ่มค่าคงที่;"`
	NameTH    string `json:"name_th" gorm:"type:varchar(100);comment:ชื่อค่าคงที่ TH;"`
	NameEN    string `json:"name_en" gorm:"type:varchar(100);comment:ชื่อค่าคงที่ EN;"`
	RefValue1 string `json:"ref_value1" gorm:"type:varchar(100);comment:ค่าอ้างอิง 1;"`
	RefValue2 string `json:"ref_value2" gorm:"type:varchar(250);comment:ค่าอ้างอิง 2;"`
	RefValue3 string `json:"ref_value3" gorm:"type:text;comment:ค่าอ้างอิง 3;"`
	Sort      int    `json:"sort" gorm:"type:int;comment:ลำดับ;"`
}
