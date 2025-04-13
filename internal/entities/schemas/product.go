package schemas

type ProductReq struct {
	ProductId   int    `json:"product_id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryID  int    `json:"category_id"`
}
