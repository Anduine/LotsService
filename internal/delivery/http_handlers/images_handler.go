package http_handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func ServeCarImage(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
  filename := vars["filename"]

	filePath := fmt.Sprintf("internal/storage/cars/%s",filename)

	//log.Printf("Запрашиваем файл: %s", filePath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(writer, "Файл зображення не знайдено", http.StatusNotFound)
		return
	}

	ext := filepath.Ext(filePath)
	if ext != ".webp" && ext != ".jpg" && ext != ".png" {
		http.Error(writer, "Доступ к зображенню заборонено", http.StatusForbidden)
		return
	}

	http.ServeFile(writer, req, filePath)
}
