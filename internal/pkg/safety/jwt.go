package safety

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConst struct {
	UserId   string `json:"user_id"`  //ผู้ใช้งาน
	Name     string `json:"name"`     //ชื่อ
	SurName  string `json:"sur_name"` //นามสกุล
	Email    string `json:"email"`    //อีเมล
	Level    string `json:"level"`    //เลเวล
	SafetyId string `json:"safety_id"`
}

type Claims struct {
	UserId   string `json:"user_id"`   //ผู้ใช้งาน
	Name     string `json:"name"`      //ชื่อ
	SurName  string `json:"sur_name"`  //นามสกุล
	Email    string `json:"email"`     //อีเมล
	Level    string `json:"level"`     //เลเวล
	SafetyId string `json:"safety_id"` //รหัสความปลอดภัย
	jwt.RegisteredClaims
}

// ฟังก์ชันสร้าง JWT
func GenerateJWT(jwtSecret string, req *JwtConst) (string, error) {
	// กำหนด claims
	claims := &Claims{
		UserId:   req.UserId,
		Name:     req.Name,
		SurName:  req.SurName,
		Email:    req.Email,
		Level:    req.Level,
		SafetyId: req.SafetyId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)), // อายุ 6 ชั่วโมง
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ArcZed", //ระบุแหล่งที่มาหรือผู้สร้างของ JWT นั้น ซึ่งมักจะเป็นชื่อของบริการ, ระบบ, หรือเซิร์ฟเวอร์ที่ออก token นี้
		},
	}

	// สร้าง token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// เซ็น token ด้วย secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT ตรวจสอบความถูกต้องของ JWT
func VerifyJWT(jwtSecret string, tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าอัลกอริธึมที่ใช้ในการเซ็นคือ HS256 หรือไม่
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// ตรวจสอบว่า token ถูกยืนยันและไม่หมดอายุ
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
