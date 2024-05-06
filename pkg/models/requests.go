package models

//Все временные struct для обработки request данных юзера. В базе данных они не сохраняются

type OrderRequest struct {
	UserID    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items"`
}

type CartItem struct {
	DishID   uint `json:"dish_id"`
	Quantity int  `json:"quantity"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AnalyticsData struct {
	TotalOrders        int64
	TotalRevenue       float64
	AverageBill        float64
	TotalUsers         int64
	MostPopularDish    DishStat
	LeastPopularDishes []DishStat
}

type DishStat struct {
	DishID   uint
	DishName string
	Count    int
}
