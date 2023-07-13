package infrastructuretype

import (
	"database/sql"
	"mistar-be-go/internal/config"
	"mistar-be-go/internal/store"
	"net/http"
)

type InfrastructureTypeHandler interface {
	GetInfrastructureTypeList(w http.ResponseWriter, r *http.Request)
	GetInfrastructureSubTypeList(w http.ResponseWriter, r *http.Request)
}

type infrastructureTypeHandler struct {
	db                      *sql.DB
	infrastructureTypeStore store.InfrastructureType
	apiCfg                  config.API
}

func NewInfrastructureTypeHandler(db *sql.DB, infrastructureTypeStore store.InfrastructureType, apiCfg config.API) InfrastructureTypeHandler {
	return &infrastructureTypeHandler{
		db:                      db,
		infrastructureTypeStore: infrastructureTypeStore,
		apiCfg:                  apiCfg,
	}
}
