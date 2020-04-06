package validates

type CreateOrderRequest struct {
	Title string `json:"username" validates:"required,gte=2,lte=50" comment:"订单标题"`
	Username string `json:"username" validates:"required,gte=2,lte=50" comment:"用户名"`
	Price float32 `json:"price" validates:"required,gte=2,lte=50" comment:"价格"`
	Description string `json:"description " validates:"required,gte=2,lte=50" comment:"描述"`
	Status string `json:"status" validates:"required,gte=2,lte=50" comment:"订单状态"`
	ImagePath string `json:"imagepath" validates:"required,gte=2,lte=50" comment:"订单图片"`
}

