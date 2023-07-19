package user

import (
	"database/sql"
	"mistar-be-go/internal/config"
	"mistar-be-go/internal/store"
	"net/http"
)

type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	db        *sql.DB
	userStore store.User
	apiCfg    config.API
}

func NewUserHandler(db *sql.DB, userStore store.User, apiCfg config.API) UserHandler {
	return &userHandler{
		db:        db,
		userStore: userStore,
		apiCfg:    apiCfg,
	}
}
