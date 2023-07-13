package infrastructure

import (
	"database/sql"
	"mistar-be-go/internal/config"
	"mistar-be-go/internal/store"
	"net/http"
)

type InfrastructureHandler interface {
	CreateInfrastructure(w http.ResponseWriter, r *http.Request)
	GetAllInfrastructureList(w http.ResponseWriter, r *http.Request)
}

type infrastructureHandler struct {
	db                      *sql.DB
	infrastructureStore     store.Infrastructure
	infrastructureTypeStore store.InfrastructureType
	apiCfg                  config.API
}

func NewInfrastructureHandler(db *sql.DB, infrastructureStore store.Infrastructure, infrastructureTypeStore store.InfrastructureType, apiCfg config.API) InfrastructureHandler {
	return &infrastructureHandler{
		db:                      db,
		infrastructureStore:     infrastructureStore,
		infrastructureTypeStore: infrastructureTypeStore,
		apiCfg:                  apiCfg,
	}
}
