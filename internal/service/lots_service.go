package service

import (
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/internal/domain"

	"github.com/google/uuid"
)

type LotsService struct {
	repo domain.LotsRepository
}

func NewLotsService(repo domain.LotsRepository) *LotsService {
	return &LotsService{repo: repo}
}

func (s *LotsService) SaveImages(files []*multipart.FileHeader) ([]string, error) {
	var savedPaths []string

	for _, fileHeader := range files {
		imgFile, err := fileHeader.Open()
		if err != nil {
			log.Println("Помилка відкриття файлу:", err)
			return nil, err
		}
		defer imgFile.Close()

		filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
		path := filepath.Join("internal/storage/cars", filename)
		dst, err := os.Create(path)
		if err != nil {
			log.Println("Помилка збереження файлу:", err)
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, imgFile); err != nil {
			log.Println("Помилка копіювання файлу:", err)
			return nil, err
		}

		savedPaths = append(savedPaths, filename)
	}
	
	return savedPaths, nil
}

func (s *LotsService) GetLotsCount() (int, error) {
	return s.repo.GetLotsCount()
}

func (s *LotsService) GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear string) (int, error) {
	return s.repo.GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear)
}

func (s *LotsService) GetLotByID(userID, lotID int) (*domain.Lot, error) {
	return s.repo.GetLotByID(userID, lotID)
}

func (s *LotsService) GetPageLots(userID, page, limit int) (*[]domain.Lot, error) {
	return s.repo.GetPageLots(userID, page, limit)
}

func (s *LotsService) GetLotsByParams(userID int, brand, model, minPrice, maxPrice, minYear, maxYear string, page, limit int) (*[]domain.Lot, error) {
	return s.repo.GetLotsByParams(userID, brand, model, minPrice, maxPrice, minYear, maxYear, page, limit)
}

func (s *LotsService) GetBrands() (*[]domain.Brand, error) {
	return s.repo.GetBrands()
}

func (s *LotsService) GetModels(brandName string) (*[]domain.Model, error) {
	return s.repo.GetModels(brandName)
}

func (s *LotsService) GetUserPostedLots(userID int) (*[]domain.Lot, error) {
	return s.repo.GetUserPostedLots(userID)
}

func (s *LotsService) GetUserLikedLots(userID int) (*[]domain.Lot, error) {
	return s.repo.GetUserLikedLots(userID)
}

func (s *LotsService) CreateLot(ctx context.Context, lot *domain.Lot) error {
	return s.repo.CreateLot(ctx, lot)
}

func (s *LotsService) UpdateLot(ctx context.Context, lot *domain.Lot) error {
	return s.repo.UpdateLot(ctx, lot)
}

func (s *LotsService) DeleteLot(ctx context.Context, lotID, userID int) error {
	lot, err := s.repo.GetLotByID(userID, lotID)
	if err != nil {
		return err
	}
	if lot.SellerID != userID {
		return errors.New("not authorized")
	}

	for _, img := range lot.Images {
		os.Remove("internal/storage/cars/" + img)
	}
	
	return s.repo.DeleteLot(ctx, lotID)
}

func (s *LotsService) LikeLot(userID, lotID int) error {
	return s.repo.LikeLot(userID, lotID)
}

func (s *LotsService) UnlikeLot(userID, lotID int) error {
	return s.repo.UnlikeLot(userID, lotID)
}

func (s *LotsService) BuyLot(userID, lotID int) error {
	lot, err := s.repo.GetLotByID(userID, lotID)
	if err != nil {
		return err
	}
	if lot.SaleStatus == "Продано" {
		return errors.New("лот уже продано")
	}

	return s.repo.MarkLotAsSold(lotID)
}
