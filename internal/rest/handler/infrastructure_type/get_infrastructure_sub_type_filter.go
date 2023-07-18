package infrastructuretype

import (
	"log"
	"net/http"

	"mistar-be-go/internal/rest/response"

	apierror "mistar-be-go/internal/rest/error"
)

type GetInfrastructureSubTypeFilterResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (handler *infrastructureTypeHandler) GetInfrastructureSubTypeFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetInfrastructureSubTypeFilterDataRequest{
		infrastructureSubTypeID: r.URL.Query().Get("infrastructure_type_id"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	infrastructureTypeList, err := handler.infrastructureTypeStore.FindAllInfrastructureTypePagination(ctx, 0, 50, 0, "", "")
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	itemsCount := len(infrastructureTypeList)

	resp := make([]GetInfrastructureSubTypeFilterResponse, itemsCount+1)
	resp[0] = GetInfrastructureSubTypeFilterResponse{
		ID:   0,
		Name: "Semua",
	}
	for idx, item := range infrastructureTypeList {

		resp[idx+1] = GetInfrastructureSubTypeFilterResponse{
			ID:   item.SubType.ID,
			Name: item.SubType.SubTypeName,
		}
	}

	response.Respond(w, http.StatusOK, resp)
}
