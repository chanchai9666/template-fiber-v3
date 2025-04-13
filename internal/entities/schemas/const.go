package schemas

import "template-fiber-v3/internal/entities/models"

type ConfigConstant struct {
	ConstId   string `json:"const_id" query:"const_id"`     //ไอดีค่าคงที่
	GroupId   string `json:"group_id" query:"group_id"`     //ไอดีกลุ่มค่าคงที่
	NameTH    string `json:"name_th" query:"name_th"`       //ชื่อค่าคงที่ TH
	NameEN    string `json:"name_en" query:"name_en"`       //ชื่อค่าคงที่ EN
	RefValue1 string `json:"ref_value1" query:"ref_value1"` //ค่าอ้างอิง 1
	RefValue2 string `json:"ref_value2" query:"ref_value2"` //ค่าอ้างอิง 2
	RefValue3 string `json:"ref_value3" query:"ref_value3"` //ค่าอ้างอิง 3
	IsActive  string `json:"is_active" query:"is_active"`   //สถานะใช้งาน
	Sort      int    `json:"sort" query:"sort"`             //ลำดับ
}

type ConfigConstantResp struct {
	models.AuditLog
	ConstId     string `json:"const_id"`                //ไอดีค่าคงที่
	GroupId     string `json:"group_id"`                //ไอดีกลุ่มค่าคงที่
	NameTH      string `json:"name_th"`                 //ชื่อค่าคงที่ TH
	NameEN      string `json:"name_en"`                 //ชื่อค่าคงที่ EN
	RefValue1   string `json:"ref_value1"`              //ค่าอ้างอิง 1
	RefValue2   string `json:"ref_value2"`              //ค่าอ้างอิง 2
	RefValue3   string `json:"ref_value3"`              //ค่าอ้างอิง 3
	IsActiveTxt string `json:"is_active_txt,omitempty"` //สถานะใช้งาน
	Sort        int    `json:"sort"`                    //ลำดับ
}
