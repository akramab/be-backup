package user

import (
	"encoding/json"
	"log"
	apierror "mistar-be-go/internal/rest/error"
	"mistar-be-go/internal/rest/response"
	"mistar-be-go/internal/store"
	"net/http"
)

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := RegisterUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	newUserData := &store.UserData{
		Name:        req.Name,
		Type:        "USER",
		Email:       req.Email,
		PhoneNumber: "+62812122912",
		Password:    req.Password,
	}

	if err := handler.userStore.Insert(ctx, newUserData); err != nil {
		log.Println("error insert new insert data: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := RegisterUserResponse{
		Message: "User created",
	}

	response.Respond(w, http.StatusCreated, resp)
}
