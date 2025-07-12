package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"server/internal/domain"
	"strings"

	"github.com/lib/pq"
)

type PostgresLotsRepo struct {
	db *sql.DB
}

func NewPostgresLotsRepo(db *sql.DB) *PostgresLotsRepo{
	return &PostgresLotsRepo{db: db}
}

func (r *PostgresLotsRepo) GetLotsCount() (int, error) {
	lotsCount := 0
	err := r.db.QueryRow("SELECT COUNT(*) FROM sell_lots").Scan(&lotsCount)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Количество лотов не найдены в базе данных ", err)
			return 0, err
		}
		log.Println("Ошибка при сканировании данных:", err)
		return 0, err
	}

	return lotsCount, nil
}

func (r *PostgresLotsRepo) GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear string) (int, error) {
	baseQuery := `
	SELECT COUNT(*)
	FROM sell_lots sl
	JOIN cars c ON sl.car_id = c.car_id
	JOIN brands b ON c.brand_id = b.brand_id
	JOIN models m ON c.model_id = m.model_id
	WHERE 1=1 
	`

	var args []any
	var conditions []string
	argCounter := 1

	addCondition := func(clause string, value any) {
		conditions = append(conditions, fmt.Sprintf("AND %s $%d", clause, argCounter))
		args = append(args, value)
		argCounter++
	}

	if brand != "" {
		addCondition("b.brand_name =", brand)
	}
	if model != "" {
		addCondition("m.model_name =", model)
	}
	if minPrice != "" {
		addCondition("sl.sale_price >=", minPrice)
	}
	if maxPrice != "" && maxPrice != "0" {
		addCondition("sl.sale_price <=", maxPrice)
	}
	if minYear != "" && minYear != "0" {
		addCondition("c.made_year >=", minYear)
	}
	if maxYear != "" && maxYear != "0" {
		addCondition("c.made_year <=", maxYear)
	}

	fullQuery := baseQuery + strings.Join(conditions, "\n")

	var lotsCount int
	err := r.db.QueryRow(fullQuery, args...).Scan(&lotsCount)
	if err != nil {
		log.Println("Ошибка при получении количества лотов:", err)
		return 0, err
	}

	return lotsCount, nil
}

func (r *PostgresLotsRepo) GetPageLots(userID, page, limit int) (*[]domain.Lot, error) {

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := `
	SELECT 
	sl.lot_id, sl.seller_id, 
  sl.postdate, sl.sale_price, sl.sale_status, sl.vin_code, 
	sl.mileage, sl.color, sl.description, sl.images_paths,
  c.car_id, c.made_year, c.engine_type, c.transmission, c.wheel_drive,
	b.brand_name, m.model_name,
	EXISTS (
		SELECT 1 FROM liked_lots ll WHERE ll.user_id = $3 AND ll.lot_id = sl.lot_id
	)
	FROM sell_lots sl
	JOIN cars c ON sl.car_id = c.car_id
	JOIN brands b ON c.brand_id = b.brand_id
	JOIN models m ON c.model_id = m.model_id
	ORDER BY postdate DESC
	LIMIT $1 OFFSET $2;
	`

	queryRows, err := r.db.Query(query, limit, offset, userID)
	if err != nil {
		log.Println("Лоты не найдены в базе данных: ", err)
		return nil, err
	}
	defer queryRows.Close()

	var lots []domain.Lot
	for queryRows.Next() {
		var lot domain.Lot
		var images pq.StringArray
		err := queryRows.Scan(
			&lot.LotID, &lot.SellerID, 
			&lot.PostDate, &lot.SalePrice, &lot.SaleStatus, &lot.VinCode, 
			&lot.Mileage, &lot.Color, &lot.Description, &images,
			&lot.Car.CarID, &lot.Car.MadeYear, &lot.Car.Engine,
			&lot.Car.Transmission, &lot.Car.WheelDrive, 
			&lot.Car.Brand, &lot.Car.Model,
			&lot.IsLiked,
		)
		if err != nil {
			log.Println("Ошибка при сканировании данных:", err)
			return nil, err
		}
		
		if images != nil {
			lot.Images = images
		} else {
			lot.Images = []string{}
		}

		lots = append(lots, lot)
	}

	return &lots, nil
}

func (r *PostgresLotsRepo) GetLotsByParams( userID int,
	brand, model, minPrice, maxPrice, minYear, maxYear string,
	page, limit int) (*[]domain.Lot, error) {

	baseQuery := `
	SELECT 
	sl.lot_id, sl.seller_id, 
	sl.postdate, sl.sale_price, sl.sale_status, sl.vin_code, 
	sl.mileage, sl.color, sl.description, sl.images_paths,
	c.car_id, c.made_year, c.engine_type, c.transmission, c.wheel_drive,
	b.brand_name, m.model_name,
	EXISTS (
		SELECT 1 FROM liked_lots ll WHERE ll.user_id = $1 AND ll.lot_id = sl.lot_id
	)
	FROM sell_lots sl
	JOIN cars c ON sl.car_id = c.car_id
	JOIN brands b ON c.brand_id = b.brand_id
	JOIN models m ON c.model_id = m.model_id
	WHERE 1=1
	`

	var args []any
	var conditions []string
	argCounter := 1

	addCondition := func(clause string, value any) {
		conditions = append(conditions, fmt.Sprintf("AND %s $%d", clause, argCounter))
		args = append(args, value)
		argCounter++
	}

	args = append(args, userID)
	argCounter++

	if brand != "" {
		addCondition("b.brand_name =", brand)
	}
	if model != "" {
		addCondition("m.model_name =", model)
	}
	if minPrice != "" {
		addCondition("sl.sale_price >=", minPrice)
	}
	if maxPrice != "" && maxPrice != "0" {
		addCondition("sl.sale_price <=", maxPrice)
	}
	if minYear != "" && minYear != "0" {
		addCondition("c.made_year >=", minYear)
	}
	if maxYear != "" && maxYear != "0" {
		addCondition("c.made_year <=", maxYear)
	}

	offset := (page - 1) * limit
	// Добавляем ORDER BY, LIMIT и OFFSET напрямую (LIMIT и OFFSET через args)
	query := baseQuery + "\n" + strings.Join(conditions, "\n") +
		fmt.Sprintf("\nORDER BY sl.postdate DESC LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, limit, offset)

	queryRows, err := r.db.Query(query, args...)
	if err != nil {
		log.Println("Лоты не найдены в базе данных: ", err)
		return nil, err
	}
	defer queryRows.Close()

	var lots []domain.Lot

	for queryRows.Next() {
		var lot domain.Lot
		var images pq.StringArray
		err := queryRows.Scan(
			&lot.LotID, &lot.SellerID, 
			&lot.PostDate, &lot.SalePrice, &lot.SaleStatus, &lot.VinCode, 
			&lot.Mileage, &lot.Color, &lot.Description, &images,
			&lot.Car.CarID, &lot.Car.MadeYear, &lot.Car.Engine,
			&lot.Car.Transmission, &lot.Car.WheelDrive, 
			&lot.Car.Brand, &lot.Car.Model,
			&lot.IsLiked,
		)
		if err != nil {
			log.Println("Ошибка при сканировании:", err)
			continue
		}

		if images != nil {
		lot.Images = images
		} else {
			lot.Images = []string{}
		}

		lots = append(lots, lot)
	}

	return &lots, nil
}

func (r *PostgresLotsRepo) GetLotByID(userID, lotID int) (*domain.Lot, error) {
	query := `
	SELECT 
  sl.lot_id, sl.seller_id, 
  sl.postdate, sl.sale_price, sl.sale_status, sl.vin_code,
	sl.mileage, sl.color, sl.description, sl.images_paths,
  c.car_id, c.made_year, c.engine_type, c.transmission, c.wheel_drive, 
	b.brand_name, b.brand_id, m.model_name, m.model_id, 
	EXISTS (
		SELECT 1 FROM liked_lots ll WHERE ll.user_id = $2 AND ll.lot_id = sl.lot_id
	)
	FROM sell_lots sl
	JOIN cars c ON sl.car_id = c.car_id
	JOIN brands b ON c.brand_id = b.brand_id
	JOIN models m ON c.model_id = m.model_id
	WHERE sl.lot_id = $1;
	`

	row := r.db.QueryRow(query, lotID, userID)

	var lot domain.Lot
	var images pq.StringArray
	err := row.Scan(
		&lot.LotID, &lot.SellerID, 
		&lot.PostDate, &lot.SalePrice, &lot.SaleStatus, &lot.VinCode, 
		&lot.Mileage, &lot.Color, &lot.Description, &images,
		&lot.Car.CarID, &lot.Car.MadeYear, &lot.Car.Engine,
		&lot.Car.Transmission, &lot.Car.WheelDrive,
		&lot.Car.Brand, &lot.Car.BrandID, &lot.Car.Model, &lot.Car.ModelID,
		&lot.IsLiked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
				log.Println("Лот не найден в базе данных ", err)
				return nil, err
		}
		log.Println("Ошибка при сканировании данных:", err)
		return nil, err
  }

	if images != nil {
		lot.Images = images
	} else {
		lot.Images = []string{}
	}

	// log.Println("Лот найден:", lot)

  return &lot, nil
}

func (r *PostgresLotsRepo) GetBrands() (*[]domain.Brand, error) {
	query := `SELECT brand_id, brand_name FROM brands`

	queryRows, err := r.db.Query(query)
	if err != nil {
		log.Println("Brands does not exist", err)
		return nil, err
	}
	defer queryRows.Close()

	var brands []domain.Brand

	for queryRows.Next() {
		var brand domain.Brand
		err := queryRows.Scan(
			&brand.BrandID, &brand.BrandName,
		)
		if err != nil {
			log.Println("Error while scanning:", err)
			return nil, err
		}

		brands = append(brands, brand)
	}

	return &brands, err
}

func (r *PostgresLotsRepo) GetModels(brandName string) (*[]domain.Model, error) {
	query := `
	SELECT m.model_id, m.brand_id, m.model_name
  FROM models m
  JOIN brands b ON m.brand_id = b.brand_id
	`

	if brandName != "" && brandName != " " {
		query += ` WHERE b.brand_name ILIKE $1`
	}

	queryRows, err := r.db.Query(query, brandName)
	if err != nil {
		log.Println("Модели не найдены в базе данных ", err)
		return nil, err
	}
	defer queryRows.Close()

	var models []domain.Model

	for queryRows.Next() {
		var model domain.Model
		err := queryRows.Scan(
			&model.ModelID, &model.BrandID, &model.ModelName,
		)
		if err != nil {
			log.Println("Ошибка при сканировании данных:", err)
			return nil, err
		}

		models = append(models, model)
	}

	return &models, err
}

func (r *PostgresLotsRepo) GetUserPostedLots(userID int) (*[]domain.Lot, error) {
	query := `
	SELECT sl.lot_id, sl.seller_id, 
  sl.postdate, sl.sale_price, sl.sale_status, sl.vin_code, 
	sl.mileage, sl.color, sl.description, sl.images_paths,
  c.car_id, c.made_year, c.engine_type, c.transmission, c.wheel_drive,
	b.brand_name, m.model_name
	FROM sell_lots sl
	JOIN cars c ON sl.car_id = c.car_id
	JOIN brands b ON c.brand_id = b.brand_id
	JOIN models m ON c.model_id = m.model_id
	WHERE seller_id = $1`

	// strUserID, _ := strconv.Atoi(userID)
	queryRows, err := r.db.Query(query, userID)
	if err != nil {
		log.Println("Лоты данного пользователя не найдены в базе данных ", err)
		return nil, err
	}
	defer queryRows.Close()

	var lots []domain.Lot
	
	for queryRows.Next() {
		var lot domain.Lot
		var images pq.StringArray
		err := queryRows.Scan(
			&lot.LotID, &lot.SellerID, 
			&lot.PostDate, &lot.SalePrice, &lot.SaleStatus, &lot.VinCode, 
			&lot.Mileage, &lot.Color, &lot.Description, &images,
			&lot.Car.CarID, &lot.Car.MadeYear, &lot.Car.Engine,
			&lot.Car.Transmission, &lot.Car.WheelDrive, 
			&lot.Car.Brand, &lot.Car.Model,
		)
		if err != nil {
			log.Println("Ошибка при сканировании:", err)
			continue
		}

		if images != nil {
			lot.Images = images
		} else {
			lot.Images = []string{}
		}

		lots = append(lots, lot)
	}

	return &lots, nil
}

func (r *PostgresLotsRepo) GetUserLikedLots(userID int) (*[]domain.Lot, error) {
	query := `
	SELECT sl.lot_id, sl.seller_id, 
  sl.postdate, sl.sale_price, sl.sale_status, sl.vin_code, 
	sl.mileage, sl.color, sl.description, sl.images_paths,
  c.car_id, c.made_year, c.engine_type, c.transmission, c.wheel_drive,
	b.brand_name, m.model_name
	FROM sell_lots sl
	JOIN liked_lots ll ON sl.lot_id = ll.lot_id
	JOIN cars c ON sl.car_id = c.car_id
	JOIN brands b ON c.brand_id = b.brand_id
	JOIN models m ON c.model_id = m.model_id
	WHERE ll.user_id = $1;`

	// strUserID, _ := strconv.Atoi(userID)
	queryRows, err := r.db.Query(query, userID)
	if err != nil {
		log.Println("Лоты лайкнутые пользователем не найдены в базе данных ", err)
		return nil, err
	}
	defer queryRows.Close()

	var lots []domain.Lot
	
	for queryRows.Next() {
		var lot domain.Lot
		var images pq.StringArray
		err := queryRows.Scan(
			&lot.LotID, &lot.SellerID, 
			&lot.PostDate, &lot.SalePrice, &lot.SaleStatus, &lot.VinCode, 
			&lot.Mileage, &lot.Color, &lot.Description, &images,
			&lot.Car.CarID, &lot.Car.MadeYear, &lot.Car.Engine,
			&lot.Car.Transmission, &lot.Car.WheelDrive, 
			&lot.Car.Brand, &lot.Car.Model,
		)
		if err != nil {
			log.Println("Ошибка при сканировании:", err)
			continue
		}

		if images != nil {
			lot.Images = images
		} else {
			lot.Images = []string{}
		}

		lots = append(lots, lot)
	}

	return &lots, nil
}

func (r *PostgresLotsRepo) CreateLot(ctx context.Context, lot *domain.Lot) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var brandID int
	err = tx.QueryRowContext(ctx, "SELECT brand_id FROM brands WHERE brand_name = $1", lot.Car.Brand).Scan(&brandID)
	if err != nil {
		log.Println("Бренд не знайдено: ", err)
		return err
	}

	var modelID int
	err = tx.QueryRowContext(ctx, `
		SELECT model_id FROM models 
		WHERE model_name = $1 AND brand_id = $2
	`, lot.Car.Model, brandID).Scan(&modelID)
	if err != nil {
		log.Println("Модель не знайдено: ", err)
		return err
	}

	var carID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO cars (made_year, engine_type, transmission, wheel_drive, brand_id, model_id, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING car_id
	`, lot.Car.MadeYear, lot.Car.Engine, lot.Car.Transmission, lot.Car.WheelDrive, brandID, modelID, lot.Description).Scan(&carID)
	if err != nil {
		log.Println("Помилка у додаванні машини: ", err)
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO sell_lots (
			seller_id, car_id, postdate, sale_price, sale_status,
			vin_code, mileage, color, description, images_paths
		)
		VALUES ($1, $2, CURRENT_DATE, $3, $4, $5, $6, $7, $8, $9)
	`,
		lot.SellerID, carID, lot.SalePrice, "Продається",
		lot.VinCode, lot.Mileage, lot.Color,
		lot.Description, pq.Array(lot.Images),
	)
	if err != nil {
		log.Println("Помилка у додаванні лота: ", err)
		log.Println("Seller Id: ", lot.SellerID)
		return err
	}

	return tx.Commit()
}

func (r *PostgresLotsRepo) UpdateLot(ctx context.Context, lot *domain.Lot) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, `
		UPDATE cars SET brand_id = $1, model_id = $2, made_year = $3, engine_type = $4, transmission = $5, wheel_drive = $6
		WHERE car_id = $7
	`, lot.Car.BrandID, lot.Car.ModelID, lot.Car.MadeYear, lot.Car.Engine, lot.Car.Transmission, lot.Car.WheelDrive, lot.Car.CarID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE sell_lots SET seller_id = $1, sale_price = $2, sale_status = $3, vin_code = $4, color = $5, mileage = $6, description = $7, images_paths = $8
		WHERE lot_id = $9
	`, lot.SellerID, lot.SalePrice, lot.SaleStatus, lot.VinCode, lot.Color, lot.Mileage, lot.Description, pq.Array(lot.Images), lot.LotID)
	
	return err
}

func (r *PostgresLotsRepo) DeleteLot(ctx context.Context, lotID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sell_lots WHERE lot_id = $1`, lotID)
	return err
}

func (r *PostgresLotsRepo) LikeLot(userID, lotID int) error {
	query := `INSERT INTO liked_lots (user_id, lot_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, userID, lotID)
	if err != nil {
		log.Println("Не вдалося додати лайк:", err)
		return err
	}

	return nil
}

func (r *PostgresLotsRepo) UnlikeLot(userID, lotID int) error {
	query := `DELETE FROM liked_lots WHERE user_id = $1 AND lot_id = $2`
	_, err := r.db.Exec(query, userID, lotID)
	if err != nil {
		log.Println("Не вдалося прибрати лайк:", err)
		return err
	}

	return nil
}

func (r *PostgresLotsRepo) MarkLotAsSold(lotID int) error {
	_, err := r.db.Exec(`UPDATE sell_lots SET sale_status = 'Продано' WHERE lot_id = $1`, lotID)
	return err
}

