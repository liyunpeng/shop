package validates

type CreateOrderRequest struct {
	Title string `json:"username" validates:"required,gte=2,lte=50" comment:"用户名"`
	Username string `json:"username" validates:"required,gte=2,lte=50" comment:"用户名"`
}

