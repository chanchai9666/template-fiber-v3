package adapter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/chanchai9666/aider"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()

type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents the structure of a JSON error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func NewSuccessMessage() *Message {
	return &Message{
		Code:    200,
		Message: "Success",
	}
}

func RespJson(c fiber.Ctx, fn interface{}, input interface{}) error {
	// ตรวจสอบและดึงข้อมูลจาก request ตาม content type
	if err := parseInputData(c, input); err != nil {
		return RenderJSON(c, err, nil)

	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validateInput(input); err != nil {
		return RenderJSON(c, err, nil)

	}

	fnType := reflect.TypeOf(fn)
	var out []reflect.Value

	// ตรวจสอบจำนวนพารามิเตอร์ที่ฟังก์ชัน fn ต้องการ
	if fnType.NumIn() == 0 {
		// ถ้าฟังก์ชันไม่ต้องการพารามิเตอร์
		out = reflect.ValueOf(fn).Call(nil)
	} else if fnType.NumIn() == 1 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 1 ตัว
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(input),
		})
	} else if fnType.NumIn() == 2 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 2 ตัว (เช่น context และ input)
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(input),
		})
	} else {
		// กรณีจำนวนพารามิเตอร์ไม่ตรง
		return RenderJSON(c, fmt.Errorf("invalid function signature"), nil)

	}

	// ตรวจสอบผลลัพธ์
	errObj := out[1].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return RenderJSON(c, errObj.(error), nil)

	}

	ResponseResult := ResponseResult{
		Code:    200,
		Message: "Success",
		Result:  out[0].Interface(),
	}
	return RenderJSON(c, nil, ResponseResult)
}

func RespJsonNoReq(c fiber.Ctx, fn interface{}) error {

	out := reflect.ValueOf(fn).Call(nil)

	// ตรวจสอบผลลัพธ์
	errObj := out[1].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return RenderJSON(c, errObj.(error), nil)

	}

	ResponseResult := ResponseResult{
		Code:    200,
		Message: "Success",
		Result:  out[0].Interface(),
	}

	return RenderJSON(c, nil, ResponseResult)
}

func RespSuccess(c fiber.Ctx, fn interface{}, input interface{}) error {
	// ตรวจสอบและดึงข้อมูลจาก request ตาม content type
	if err := parseInputData(c, input); err != nil {
		return RenderJSON(c, err, nil)

	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := validateInput(input); err != nil {
		return RenderJSON(c, err, nil)

	}

	fnType := reflect.TypeOf(fn)
	var out []reflect.Value

	// ตรวจสอบจำนวนพารามิเตอร์ที่ฟังก์ชัน fn ต้องการ
	if fnType.NumIn() == 0 {
		// ถ้าฟังก์ชันไม่ต้องการพารามิเตอร์
		out = reflect.ValueOf(fn).Call(nil)
	} else if fnType.NumIn() == 1 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 1 ตัว
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(input),
		})
	} else if fnType.NumIn() == 2 {
		// ถ้าฟังก์ชันต้องการพารามิเตอร์ 2 ตัว (เช่น context และ input)
		out = reflect.ValueOf(fn).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(input),
		})
	} else {
		// กรณีจำนวนพารามิเตอร์ไม่ตรง
		return RenderJSON(c, fmt.Errorf("invalid function signature"), nil)

	}
	// ตรวจสอบผลลัพธ์
	errObj := out[0].Interface()
	if errObj != nil {
		logrus.Errorf("call service error: %s", errObj)
		return RenderJSON(c, errObj.(error), nil)

	}

	return RenderJSON(c, nil, NewSuccessMessage())
}

// ฟังก์ชันนี้ใช้ในการเช็คและดึงข้อมูลจาก request ตาม content type
func parseInputData(c fiber.Ctx, input interface{}) error {
	Method := c.Method()
	switch {
	case strings.HasPrefix(Method, "POST"), strings.HasPrefix(Method, "DELETE"):
		// เรียกใช้ฟังก์ชันแปลง
		if err := parseRequestBody(c, &input); err != nil {
			return err
		}
	case strings.HasPrefix(Method, "GET"):
		// ถ้าต้องการจัดการกับ query parameters ด้วย
		if err := c.Bind().Query(input); err != nil {
			return err
		}
	case c.Get("Content-Type") == "multipart/form-data":
		if err := mapFormValues(c, input); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported content type: %s", Method)
	}

	// ดึง path parameters แล้วแปลงใส่ struct
	val := reflect.ValueOf(input).Elem()
	for _, name := range c.Route().Params {
		value := c.Params(name)

		field := val.FieldByName(toPascalCase(name))
		if field.IsValid() && field.CanSet() {
			field.SetString(value)
		}
	}

	return nil
}

func mapFormValues(c fiber.Ctx, input interface{}) error {
	val := reflect.ValueOf(input).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// ใช้ c.FormValue แทน c.PostForm ใน Fiber v3
		value := c.FormValue(fieldName)
		if value != "" {
			if field.Type().Kind() == reflect.Slice {
				slice := reflect.MakeSlice(field.Type(), 1, 1)
				slice.Index(0).SetString(value)
				field.Set(slice)
			} else if field.Kind() == reflect.String {
				field.SetString(value)
			} else {
				// รองรับประเภทอื่นๆ เช่น int, bool, float
				if err := setValue(field, value); err != nil {
					return fmt.Errorf("cannot set field %s: %w", fieldName, err)
				}
			}
		}
	}

	return nil
}

func parseRequestBody(c fiber.Ctx, inputStruct interface{}) error {
	if err := c.Bind().Body(inputStruct); err != nil {
		return fmt.Errorf("failed to parse request body: %v", err)
	}
	return nil
}

// ฟังก์ชันที่จัดการ error response และ success response
func RenderJSON(c fiber.Ctx, err error, successResponse interface{}) error {
	if err != nil {
		// สร้าง ErrorResponse instance
		errorResponse := ErrorResponse{
			Code:    400, // หรือใช้รหัสข้อผิดพลาดที่คุณต้องการ
			Message: err.Error(),
		}

		// เช็คว่าข้อผิดพลาดเป็นชนิด *CustomError หรือไม่
		if customErr, ok := err.(*aider.CustomError); ok {
			// เข้าถึงรหัสข้อผิดพลาดโดยตรง
			errorResponse.Code = customErr.Code
			errorResponse.Message = customErr.Message
		}

		// ตั้งค่ารหัสสถานะ HTTP
		return c.Status(400).JSON(errorResponse)

	}

	// ส่งผลลัพธ์สำเร็จเป็น JSON
	return c.Status(200).JSON(successResponse)
}

// CustomValidator คือ struct ที่ใช้เพื่อสร้าง custom validator
type CustomValidator struct {
	validator *validator.Validate
}

// ValidateDate คือ custom validation function สำหรับการตรวจสอบวันที่
func ValidateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	// ตรวจสอบว่าค่าวันที่ไม่เป็นค่าว่าง
	if dateStr == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// NewValidator คือฟังก์ชั่นที่ใช้สร้าง custom validator
func NewValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("date", ValidateDate)

	return &CustomValidator{validator: v}
}

// ฟังก์ชันนี้ใช้ในการตรวจสอบความถูกต้องของข้อมูล
func validateInput(input interface{}) error {
	// สร้าง custom validator
	validate := NewValidator()

	// ใช้ custom validator เพื่อ validate ข้อมูล
	if err := validate.validator.Struct(input); err != nil {
		// กรณีมี error ในการ validate แสดงข้อความเพิ่มเติม
		errs := err.(validator.ValidationErrors)
		errorMsg := "Invalid request data:"
		for _, e := range errs {
			errorMsg += fmt.Sprintf("\n- Field: %s, Type: %T, Error: %s", e.Field(), e.Value(), e.Tag())
		}
		// ใช้ gin.Error เพื่อสร้าง error ที่เข้ากับ Gin framework
		return fmt.Errorf("%s", errorMsg)
	}
	return nil
}

// ฟังก์ชันนี้จะแปลง snake_case เป็น PascalCase
func toPascalCase(snake string) string {
	// แยก string ตามขีดล่าง _
	parts := strings.Split(snake, "_")
	for i, part := range parts {
		// แปลงตัวอักษรแรกของแต่ละส่วนให้เป็นตัวพิมพ์ใหญ่
		parts[i] = strings.Title(part)
	}
	// รวมกลับมาเป็น string เดียวในรูปแบบ PascalCase
	return strings.Join(parts, "")
}

// ฟังก์ชันนี้จะแปลง snake_case เป็น camelCase
func toCamelCase(snake string) string {
	pascal := toPascalCase(snake)
	// แปลงตัวอักษรแรกให้เป็นตัวพิมพ์เล็กสำหรับ camelCase
	return string(unicode.ToLower(rune(pascal[0]))) + pascal[1:]
}

func setValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("unsupported kind: %s", field.Kind())
	}
	return nil
}
