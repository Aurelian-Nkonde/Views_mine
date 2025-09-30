package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"thousand.views_mine/cmd/helpers"
	"thousand.views_mine/internals/database/db_quaries"
)

func (h *App) Login(w http.ResponseWriter, r *http.Request) {
	type loginData struct {
		Username string
		Password string
	}

	var data loginData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("error getting data")
		http.Error(w, "error getting data", http.StatusBadRequest)
		return
	}

	acc, err := h.Quaries.GetUserByAccount(h.Ctx, data.Username)
	if err != nil {
		log.Println("account is not found")
		http.NotFound(w, r)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(data.Password))
	if err != nil {
		log.Println("password err")
		http.Error(w, "password err", http.StatusBadRequest)
		return
	}

	type accountinfo struct {
		Id        string
		Username  string
		Email     string
		Token     string
		Verified  bool
		CreatedAt pgtype.Timestamp
	}
	_, tns, err := h.Token.Encode(map[string]interface{}{
		"user_id": acc.AccountID,
	})

	if err != nil {
		log.Println("error generating jwt")
		http.Error(w, "error generating jwt", http.StatusBadRequest)
		return
	}

	updatedData := accountinfo{
		Id:        acc.AccountID.String(),
		Username:  acc.Username,
		Email:     acc.Email,
		Token:     tns,
		Verified:  acc.EmailVerified,
		CreatedAt: acc.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(updatedData); err != nil {
		log.Println("error encoding data")
		w.Write([]byte("error!"))
		return
	}
}

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), 14)
	if err != nil {
		log.Println("error hashing the password!")
		http.Error(w, "error hashing the password", http.StatusInternalServerError)
		return
	}

	new_account, err := h.Quaries.CreateAccount(h.Ctx, db_quaries.CreateAccountParams{
		Username:  account.Username,
		Email:     account.Email,
		Password:  string(hashedPassword),
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
