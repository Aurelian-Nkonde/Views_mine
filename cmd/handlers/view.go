package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"thousand.views_mine/cmd/helpers"
	"thousand.views_mine/internals/database/db_quaries"
)

func (h *App) GetAllViews(w http.ResponseWriter, r *http.Request) {
	views, err := h.Quaries.GetAllViews(h.Ctx)
	if err != nil {
		log.Println("Error getting views")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error getting views"))
		return
	}
	if views == nil {
		views = []db_quaries.View{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(views); err != nil {
		log.Println("Error encoding json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *App) CreateView(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data db_quaries.View
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("error decoding data")
		http.Error(w, "error decoding data", http.StatusBadRequest)
		return
	}

	// simple validations!
	if data.Title == "" || data.Paragraph == "" || data.UserID == uuid.Nil || data.ViewID == uuid.Nil {
		log.Println("data fields must not be empty")
		http.Error(w, "data fields must not be empty", http.StatusBadRequest)
		return
	}

	rnd, err := uuid.NewRandom()
	if err != nil {
		log.Println("error generating uuid")
		http.Error(w, "error generating uuid", http.StatusInternalServerError)
		return
	}

	new_view, err := h.Quaries.CreateView(h.Ctx, db_quaries.CreateViewParams{
		Title:     data.Title,
		Paragraph: data.Paragraph,
		UserID:    data.UserID,
		Public:    data.Public,
		ViewID:    rnd,
	})
	if err != nil {
		log.Println("error creating a new view")
		http.Error(w, "error creating a new view", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(new_view); err != nil {
		log.Println("error encoding data")
		w.Write([]byte("error encoding data"))
	}
}

func (h *App) GetView(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	helpers.CheckIdIsNotEmpty(id, w)

	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing id")
		http.Error(w, "Error parsing id", http.StatusBadRequest)
		return
	}
	view, err := h.Quaries.GetView(h.Ctx, parsedId)
	if err != nil {
		log.Println("Error getting a view")
		http.Error(w, "error getting a view", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(view); err != nil {
		log.Println("error encoding json")
		w.Write([]byte("error encoding json"))
		return
	}
}

func (h *App) GetUserViews(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	helpers.CheckIdIsNotEmpty(id, w)

	parseduuid, err := uuid.Parse(id)
	if err != nil {
		log.Println("error parsing the id")
		http.Error(w, "error parsing the id", http.StatusBadRequest)
		return
	}

	views, err := h.Quaries.GetUserViews(h.Ctx, parseduuid)
	if err != nil {
		log.Println("error finding views")
		http.Error(w, "error finding views", http.StatusNotFound)
		return
	}

	if views == nil {
		views = []db_quaries.View{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(views); err != nil {
		log.Println("error encoding json")
		w.Write([]byte("error encoding json"))
		return
	}
}

func (h *App) GetAllPublicViews(w http.ResponseWriter, r *http.Request) {
	views, err := h.Quaries.GetAllPublicViews(h.Ctx)
	if err != nil {
		log.Println("error getting public views")
		http.Error(w, "error getting public views", http.StatusNotFound)
		return
	}

	if views == nil {
		views = []db_quaries.View{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(views); err != nil {
		log.Println("error encoding json")
		w.Write([]byte("error encoding json"))
	}
}

func (h *App) DeleteView(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	helpers.CheckIdIsNotEmpty(id, w)

	parsedid, err := uuid.Parse(id)
	if err != nil {
		log.Println("error parsing id")
		http.Error(w, "error parsing id", http.StatusBadRequest)
		return
	}

	err = h.Quaries.DeleteView(h.Ctx, parsedid)
	if err != nil {
		log.Println("error deleting a view")
		http.Error(w, "error deleting a view", http.StatusInternalServerError)
		return
	}

	log.Println("deleting complete")
	w.WriteHeader(http.StatusNoContent)
}
