package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Nandainthegrass/Flipzon/Authentication/services/auth"
	"github.com/Nandainthegrass/Flipzon/Authentication/types"
	"github.com/Nandainthegrass/Flipzon/Authentication/utils"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
	rdb   *redis.Client
}

func NewHandler(store types.UserStore, rdb *redis.Client) *Handler {
	return &Handler{
		store: store,
		rdb:   rdb,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
}

// Register Handler
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	//parse json into go readable format
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate user payload wrt to tags in the type
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload %v", errors))
		return
	}

	//check if user already exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if u != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User with email %v already exists", payload.Email))
		return
	}

	payload.Password, err = auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error occured while hashing password", err))
		return
	}
	newUser := &types.User{
		ID:       auth.GenerateID(),
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
		Phone:    payload.Phone,
	}
	err = h.store.CreateUser(*newUser)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("Error while trying to create user account", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"success": "Account created Successfully"})
}

// Login Handler
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

	var payload types.LoginUserPayload

	//decode r.Body into payload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Invalid payload %v", errors))
		return
	}

	//check if user exists
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if u == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("User not found"))
		return
	}

	//check if passwords match
	check := auth.VerifyPassword(payload.Password, u.Password)
	if !check {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Wrong password"))
		return
	}

	//Need to implement session logic
	SessionID := auth.GenerateID()
	auth.SetCookie(w, SessionID)
	//Now then we have to set redis
	err = h.rdb.LPush(context.Background(), SessionID, "").Err()
	if err != nil {
		fmt.Println("Error setting value:", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"success": "Login successful"})
}
