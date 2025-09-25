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

func (h *App) GetAccounts(w http.ResponseWriter, r *http.Request) {
	accountsData, err := h.Quaries.GetAccounts(h.Ctx)
	if err != nil {
		log.Println("Error getting accounts!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	if accountsData == nil {
		accountsData = []db_quaries.Account{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(accountsData); err != nil {
		log.Println("error encoding json")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *App) GetAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	helpers.CheckIdIsNotEmpty(id, w)
	userid, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing id to uuid")
		http.Error(w, "error parsing uuid", http.StatusBadRequest)
		return
	}
	account, err := h.Quaries.GetAccount(h.Ctx, userid)
	if err != nil {
		log.Println("Error getting an account")
		http.Error(w, "Error getting an account", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		log.Println("Error encoding json")
		w.Write([]byte("Error encoding json"))
		return
	}

}

func (h *App) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("verify email!"))
}

func (h *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	helpers.CheckIdIsNotEmpty(id, w)
	parsedId, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error parsing id to uuid")
		http.Error(w, "error parsing uuid", http.StatusBadRequest)
		return
	}
	err = h.Quaries.DeleteAccount(h.Ctx, parsedId)
	if err != nil {
		log.Println("Error deleting the account!")
		http.Error(w, "Error deleting the account", http.StatusInternalServerError)
		return
	}
	log.Println("Complete deleting id")
	w.WriteHeader(http.StatusNoContent)

}

func (h *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var account db_quaries.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		log.Println("Error decoding body data")
		http.Error(w, "Error decoding body", http.StatusBadRequest)
		return
	}

	if account.Username == "" || account.Email == "" || account.Password == "" {
		log.Println("Empty fields are not allowed!")
		http.Error(w, "Empty fields error", http.StatusBadRequest)
		return
	}

	rnd, err := uuid.NewRandom()
	if err != nil {
		log.Println("error generating a random uuid")
		http.Error(w, "Error generating a random uuid", http.StatusInternalServerError)
		return
	}
	new_account, err := h.Quaries.CreateAccount(h.Ctx, db_quaries.CreateAccountParams{
		Username:  account.Username,
		Email:     account.Email,
		Password:  account.Password,
		AccountID: rnd,
	})
	if err != nil {
		log.Println("Error creating new account!, " + err.Error())
		http.Error(w, "Error creating new account!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(new_account); err != nil {
		log.Println("Error encoding json")
		w.WriteHeader(http.StatusInternalServerError)
	}

}
