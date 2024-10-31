package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/m21power/Ecom/config"
	"github.com/m21power/Ecom/service/auth"
	"github.com/m21power/Ecom/types"
	"github.com/m21power/Ecom/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		erro := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payloads: %v", erro))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s not found", payload.Email))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check of password
	if !auth.ComparePassword(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		erro := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payloads: %v", erro))
		return
	}
	// check if the user exist
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exist", payload.Email))
		return
	}
	hashedPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// if doesn't we create
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)

}
