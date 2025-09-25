package helpers

import (
	"log"
	"net/http"
)

func CheckIdIsNotEmpty(id string, w http.ResponseWriter) {
	if id == "" {
		log.Println("parameter is empty string")
		http.Error(w, "empty parameter is not allowed!", http.StatusBadRequest)
		return
	}
}
