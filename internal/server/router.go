package server

import (
	"net/http"
	"server/internal/delivery/http_handlers"
	"server/pkg/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRouter(lotsHandler *http_handlers.LotsHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/lots/sell_lots_count", lotsHandler.GetLotsCount).Methods("GET")
	router.HandleFunc("/api/lots/sell_lots_filtered_count", lotsHandler.GetLotsByParamsCount).Methods("GET")
	router.Handle("/api/lots/sell_lots_id/{lot_id}", auth.OptionalAuthMiddleware(http.HandlerFunc(lotsHandler.GetLotByID))).Methods("GET")
	router.Handle("/api/lots/sell_lots", auth.OptionalAuthMiddleware(http.HandlerFunc(lotsHandler.GetLotsPage))).Methods("GET")
	router.Handle("/api/lots/sell_lots_filtered", auth.OptionalAuthMiddleware(http.HandlerFunc(lotsHandler.GetLotsPageByParams))).Methods("GET")

	router.HandleFunc("/api/lots/brands", lotsHandler.GetBrands).Methods("GET")
	router.HandleFunc("/api/lots/models", lotsHandler.GetModels).Methods("GET")

	router.Handle("/api/lots/user_posted_lots", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.GetUserPostedLots))).Methods("GET")
	router.Handle("/api/lots/user_liked_lots", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.GetUserLikedLots))).Methods("GET")

	router.Handle("/api/lots/create_lot", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.CreateLot))).Methods("POST")
	router.Handle("/api/lots/update_lot/{lot_id}", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.UpdateLot))).Methods("PUT")
	router.Handle("/api/lots/delete_lot/{lot_id}", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.DeleteLot))).Methods("DELETE")

	router.Handle("/api/lots/likes/{lot_id}", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.LikeLot))).Methods("POST")
	router.Handle("/api/lots/likes/{lot_id}", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.UnlikeLot))).Methods("DELETE")

	router.Handle("/api/lots/buy_lot/{lot_id}", auth.AuthMiddleware(http.HandlerFunc(lotsHandler.BuyLotHandler))).Methods("PUT")


	router.HandleFunc("/api/lots/images/{filename}", http_handlers.ServeCarImage).Methods("GET")

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	return handler
}

