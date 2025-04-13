package schemas

type AddUsers struct {
	UserId      uint64 `json:"user_id"`                  //ผู้ใช้งาน
	Email       string `json:"email" binding:"required"` //อีเมล
	Password    string `json:"password"`                 //รหัสผ่าน
	Name        string `json:"name"`                     //ชื่อ
	SurName     string `json:"sur_name"`                 //นามสกุล
	PhoneNumber string `json:"phone_number"`             //เบอร์โทร
	IdCard      string `json:"id_card"`                  //เลขบัตรประจำตัว
}

type FindUsersReq struct {
	UserId  string `json:"user_id" query:"user_id"`   //ผู้ใช้งาน
	Name    string `json:"name" query:"name"`         //ชื่อ
	SurName string `json:"sur_name" query:"sur_name"` //นามสกุล
	Email   string `json:"email" query:"email"`       //อีเมล
}

type FindUsersByEmailReq struct {
	Email string `json:"email" query:"email"` //ผู้ใช้งาน
}

type ValueReq struct {
	Value string `query:"value"` //ค่า string ใดๆ
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"  example:"admin@admin.com"` //ผู้ใช้งาน
	Password string `json:"password" binding:"required"  example:"1234"`         //รหัสผ่าน
}

type UserResp struct {
	UserId  uint64   `json:"user_id"`  //ผู้ใช้งาน
	Email   string   `json:"email"`    //อีเมล
	Name    string   `json:"name"`     //ชื่อ
	SurName string   `json:"sur_name"` //นามสกุล
	Level   []string `json:"level"`    //เลเวล
}

type LoginResp struct {
	AccessToken string   `json:"access_token"` //Token เข้าใช้งาน
	User        UserResp `json:"user"`         //ข้อมูลผู้ใช้
}

type JwtReq struct {
	UserId  string `json:"user_id"`  //ผู้ใช้งาน
	Name    string `json:"name"`     //ชื่อ
	SurName string `json:"sur_name"` //นามสกุล
	Email   string `json:"email"`    //อีเมล
	Level   string `json:"level"`    //เลเวล
}

type RefreshTokenReq struct {
	UserId uint64 `json:"user_id" binding:"required"` //ผู้ใช้งาน
	Email  string `json:"email" binding:"required"`   //อีเมล
}
