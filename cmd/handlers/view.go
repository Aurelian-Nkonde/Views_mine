package handlers

import "net/http"

func (h *App) GetAllViews(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get all views"))
}
