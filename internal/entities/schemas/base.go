package schemas

type HTTPError struct {
	Status  string
	Message string
}

// mock struct เอาไว้ให้ endpoint เรียก
type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty" query:"limit"`
	Page       int    `json:"page,omitempty" query:"page"`
	Sort       string `json:"sort,omitempty" query:"sort"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       []T    `json:"rows" swaggertype:"array,object"` //เพิ่ม swaggertype:"array,object" เพื่อให้ swagger init ได้
}
