package domain

import "context"

type Brand struct {
	BrandID   int    `json:"brand_id"`
	BrandName string `json:"brand_name"`
}

type Model struct {
	ModelID   int    `json:"model_id"`
	BrandID   int    `json:"brand_id"`
	ModelName string `json:"model_name"`
}

type Car struct {
	CarID        int    `json:"car_id"`
	BrandID			 int		`json:"brand_id"`
	ModelID      int		`json:"model_id"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	MadeYear     int    `json:"made_year"`
	Engine     	 string `json:"engine_type"`
	Transmission string `json:"transmission"`
	WheelDrive   string `json:"wheel_drive"`
}

type Lot struct {
	LotID       int      `json:"lot_id"`
	SellerID    int      `json:"seller_id"`
	Car         Car      `json:"car"`
	PostDate    string   `json:"postdate"`
	SalePrice   int      `json:"sale_price"`
	SaleStatus  string   `json:"sale_status"`
	VinCode     string   `json:"vin_code"`
	Color       string   `json:"color"`
	Mileage     int      `json:"mileage"`
	Description string   `json:"description"`
	IsLiked     bool     `json:"is_liked"`
	Images      []string `json:"images"`
}

type LotsRepository interface {
	GetLotsCount() (int, error)
	GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear string) (int, error)
	GetLotByID(userID, lotID int) (*Lot, error)
	GetPageLots(userID, page, limit int) (*[]Lot, error)
	GetLotsByParams(userID int, brand, model, minPrice, maxPrice, minYear, maxYear string, page, limit int) (*[]Lot, error)

	GetBrands() (*[]Brand, error)
	GetModels(brandName string) (*[]Model, error)

	GetUserPostedLots(userID int) (*[]Lot, error)
	GetUserLikedLots(userID int) (*[]Lot, error)

	CreateLot(ctx context.Context, lot *Lot) error
	UpdateLot(ctx context.Context, lot *Lot) error
	DeleteLot(ctx context.Context, lotID int) error
	
	LikeLot(userID, lotID int) error
	UnlikeLot(userID, lotID int) error

	MarkLotAsSold(lotID int) error
}
