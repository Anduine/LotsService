package http_handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/internal/domain"
	"server/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type LotsHandler struct {
	service *service.LotsService
}

func NewLotsHandler(service *service.LotsService) *LotsHandler {
	return &LotsHandler{service: service}
}

func (h *LotsHandler) GetLotsCount(w http.ResponseWriter, r *http.Request) {

	lotsCount, err := h.service.GetLotsCount()
	if err != nil {
		http.Error(w, "Лоти не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lotsCount); err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
	}
}

func (h *LotsHandler) GetLotsByParamsCount(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	brand := params.Get("brand")
	model := params.Get("model")
	minPrice := params.Get("min_price")
	maxPrice := params.Get("max_price")
	minYear := params.Get("min_year")
	maxYear := params.Get("max_year")

	lotsCount, err := h.service.GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear)
	if err != nil {
		http.Error(w, "Лоти за вказаними параметрами не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lotsCount); err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
	}
}

func (h *LotsHandler) GetLotByID(w http.ResponseWriter, r *http.Request) {
	var userID int
	if uidRaw := r.Context().Value("user_id"); uidRaw != nil {
			if uid, ok := uidRaw.(int); ok {
					userID = uid
			}
	}


	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		http.Error(w, "Некоректний ID лота", http.StatusBadRequest)
		return
	}

	lot, err := h.service.GetLotByID(userID, lotID)
	if err != nil {
		http.Error(w, "Лот не знайдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lot); err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
	}
}

func (h *LotsHandler) GetLotsPage(w http.ResponseWriter, r *http.Request) {
	var userID int
	if uidRaw := r.Context().Value("user_id"); uidRaw != nil {
			if uid, ok := uidRaw.(int); ok {
					userID = uid
			}
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	lots, err := h.service.GetPageLots(userID, page, limit)
	if err != nil {
		http.Error(w, "Лоти не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lots); err != nil {
		log.Println("Ошибка при кодировании JSON: ", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
	}	
}

func (h *LotsHandler) GetLotsPageByParams(w http.ResponseWriter, r *http.Request) {
	var userID int
	if uidRaw := r.Context().Value("user_id"); uidRaw != nil {
			if uid, ok := uidRaw.(int); ok {
					userID = uid
			}
	}

	params := r.URL.Query()

	brand := params.Get("brand")
	model := params.Get("model")
	minPrice := params.Get("min_price")
	maxPrice := params.Get("max_price")
	minYear := params.Get("min_year")
	maxYear := params.Get("max_year")

	pageStr := params.Get("page")
	limitStr := params.Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	lots, err := h.service.GetLotsByParams(userID, brand, model, minPrice, maxPrice, minYear, maxYear, page, limit)
	if len(*lots) == 0 || err != nil {
		http.Error(w, "Лоти за вказаними параметрами не знайдені", http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lots); err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
	}
}

func (h *LotsHandler) GetBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := h.service.GetBrands()
	if err != nil {
		http.Error(w, "Бренди не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(brands); err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
	}
}

func (h *LotsHandler) GetModels(w http.ResponseWriter, r *http.Request) {
	brandName := r.URL.Query().Get("brand")

	models, err := h.service.GetModels(brandName)
	if err != nil {
		http.Error(w, "Моделі не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

func (h *LotsHandler) GetUserPostedLots(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	postedLots, err := h.service.GetUserPostedLots(userID)
	if err != nil {
		http.Error(w, "Опубликовані користувачем лоти не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postedLots)
}

func (h *LotsHandler) GetUserLikedLots(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	LikedLots, err := h.service.GetUserLikedLots(userID)
	if err != nil {
		http.Error(w, "Лайкнуті користувачем лоти не знайдені", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LikedLots)
}

func (h *LotsHandler) CreateLot(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Помилка парсингу форми", http.StatusBadRequest)
		return
	}

	sellerID := r.Context().Value("user_id").(int)

	car := domain.Car{
		Brand:        r.FormValue("brand"),
		Model:        r.FormValue("model"),
		Engine:       r.FormValue("engine_type"),
		Transmission: r.FormValue("transmission"),
		WheelDrive:   r.FormValue("wheel_drive"),
	}
	lot := domain.Lot{
		SellerID:    sellerID,
		Car:         car,
		Color:       r.FormValue("color"),
		VinCode:     r.FormValue("vin_code"),
		Description: r.FormValue("description"),
	}

	var errConv error
	if lot.SalePrice, errConv = strconv.Atoi(r.FormValue("sale_price")); errConv != nil {
		log.Println("Невірна ціна:", errConv)
		http.Error(w, "Невірна ціна", http.StatusBadRequest)
		return
	}
	
	if lot.Mileage, errConv = strconv.Atoi(r.FormValue("mileage")); errConv != nil {
		log.Println("Невірний пробіг:", errConv)
		http.Error(w, "Невірний пробіг", http.StatusBadRequest)
		return
	}
	
	if lot.Car.MadeYear, errConv = strconv.Atoi(r.FormValue("made_year")); errConv != nil {
		log.Println("Невірна дата виробництва: ", errConv)
		http.Error(w, "Невірна дата виробництва", http.StatusBadRequest)
		return
	}

	// for _, fh := range r.MultipartForm.File["new_images"] {
  //   log.Println("File header:", fh.Filename)
	// }

	files := r.MultipartForm.File["new_images"]
	imageFilenames, err := h.service.SaveImages(files)
	if err != nil {
		log.Println("Помилка збереження зображень:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	lot.Images = imageFilenames

	if err := h.service.CreateLot(r.Context(), &lot); err != nil {
		log.Println("Помилка збереження лота:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LotsHandler) UpdateLot(w http.ResponseWriter, r *http.Request) {
	sellerID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		log.Println("Некоректний ID лота:", err)
		http.Error(w, "Некоректний ID лота", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Помилка парсингу форми", http.StatusBadRequest)
		return
	}

	formIntValue := func(key string) int {
		formValue := r.FormValue(key) 
		result, err := strconv.Atoi(formValue) 
		if err != nil {
			log.Println("Некоректне значення для", key, ":", formValue, "-", err)
			http.Error(w, "Некоректне значення для "+key, http.StatusBadRequest)
			return 0
		}
		return result 
	}

	if formIntValue("seller_id") != sellerID {
		log.Println("Не співпадає id користувача з форми:", formIntValue("seller_id"), "з JWT-токена:", sellerID)
		http.Error(w, "Не авторизовано", http.StatusUnauthorized)
		return
	}

	car := domain.Car{
		CarID:        formIntValue("car_id"),
		BrandID:      formIntValue("brand_id"),
		ModelID:      formIntValue("model_id"),
		Brand:        r.FormValue("brand"),
		Model:        r.FormValue("model"),
		MadeYear:     formIntValue("made_year"),
		Engine:       r.FormValue("engine_type"),
		Transmission: r.FormValue("transmission"),
		WheelDrive:   r.FormValue("wheel_drive"),
	}

	lot := domain.Lot{
		LotID:       lotID,	
		SellerID:    sellerID,
		Car:         car,
		Color:       r.FormValue("color"),
		VinCode:     r.FormValue("vin_code"),
		Description: r.FormValue("description"),
		SaleStatus:  r.FormValue("sale_status"),
		SalePrice:   formIntValue("sale_price"),
		Mileage:     formIntValue("mileage"),
	}

	deleteImages := r.MultipartForm.Value["delete_images"]
	newFiles := r.MultipartForm.File["new_images"]
	oldImagesStr := r.FormValue("old_images")
	var oldImages []string
	_ = json.Unmarshal([]byte(oldImagesStr), &oldImages)

	for _, img := range deleteImages {
		_ = os.Remove("internal/storage/cars/" + img)
	}

	newImageNames, err := h.service.SaveImages(newFiles)
	if err != nil {
		log.Println("Помилка збереження зображень:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	lot.Images = append(oldImages, newImageNames...)

	err = h.service.UpdateLot(r.Context(), &lot)
	if err != nil {
		log.Println("Помилка збереження лота:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LotsHandler) DeleteLot(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		log.Println("Некоректний ID лота:", err)
		http.Error(w, "Невірний ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteLot(r.Context(), lotID, userID)
	if err != nil {
		log.Println("Помилка видалення лота:", err)
		http.Error(w, "Помилка на сервері", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LotsHandler) LikeLot(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		log.Println("Некоректний ID лота:", err)
		http.Error(w, "Некоректний ID лота", http.StatusBadRequest)
		return
	}

	err = h.service.LikeLot(userID, lotID)
	if err != nil {
		http.Error(w, "Не вдалося додати лайк", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LotsHandler) UnlikeLot(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		log.Println("Некоректний ID лота:", err)
		http.Error(w, "Некоректний ID лота", http.StatusBadRequest)
		return
	}

	err = h.service.UnlikeLot(userID, lotID)
	if err != nil {
		http.Error(w, "Не вдалося прибрати лайк", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *LotsHandler) BuyLotHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	vars := mux.Vars(r)
	lotID, err := strconv.Atoi(vars["lot_id"])
	if err != nil {
		log.Println("Некоректний ID лота:", err)
		http.Error(w, "Некоректний ID лота", http.StatusBadRequest)
		return
	}

	err = h.service.BuyLot(userID, lotID)
	if err != nil {
		log.Println("Помилка купівлі лота:", err)
		http.Error(w, "Неможливо купити лот: ", http.StatusBadRequest)
		return
	}

	//log.Println("Лот куплений користувачем с ID:", userID, "лот ID:", lotID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Лот куплено"))
}