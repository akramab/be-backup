package infrastructuretype

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"mistar-be-go/internal/rest/response"

	apierror "mistar-be-go/internal/rest/error"
)

type GetInfrastructureSubTypeFilterDataRequest struct {
	infrastructureSubTypeID     string
	infrastructureSubTypeIDList []int
}

func (r *GetInfrastructureSubTypeFilterDataRequest) validate() *apierror.FieldError {
	fieldErr := apierror.NewFieldError()

	r.infrastructureSubTypeID = strings.TrimSpace(r.infrastructureSubTypeID)

	infrastructureSubTypeIDListTemp := []int{}

	infrastructureSubTypeIDListString := strings.Split(r.infrastructureSubTypeID, ",")

	for _, subTypeID := range infrastructureSubTypeIDListString {
		subTypeIDInt, err := strconv.Atoi(subTypeID)
		if err != nil && strings.TrimSpace(subTypeID) != "" {
			fieldErr = fieldErr.WithField("infrastructure_id", "invalid infrastructure id")
		} else if strings.TrimSpace(subTypeID) == "" {
			continue
		} else if subTypeIDInt == 0 {
			r.infrastructureSubTypeIDList = []int{}
			return nil
		} else {
			infrastructureSubTypeIDListTemp = append(infrastructureSubTypeIDListTemp, subTypeIDInt)
		}
	}

	r.infrastructureSubTypeIDList = infrastructureSubTypeIDListTemp

	if len(fieldErr.Fields) != 0 {
		return &fieldErr
	}

	return nil
}

type GetInfrastructureSubTypeFilterDataResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	InfrastructureCount int    `json:"infrastructure_count"`
}

func (handler *infrastructureTypeHandler) GetInfrastructureSubTypeFilterData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetInfrastructureSubTypeFilterDataRequest{
		infrastructureSubTypeID: r.URL.Query().Get("infrastructure_type_id"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	infrastructureSubTypeFilterDataList, err := handler.infrastructureTypeStore.FindInfrastructureSubTypeFilterData(ctx, req.infrastructureSubTypeIDList)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	itemsCount := len(infrastructureSubTypeFilterDataList)

	resp := make([]GetInfrastructureSubTypeFilterDataResponse, itemsCount)
	for idx, item := range infrastructureSubTypeFilterDataList {
		resp[idx] = GetInfrastructureSubTypeFilterDataResponse{
			ID:                  item.ID,
			Name:                fmt.Sprintf("Infrastruktur %s", item.SubTypeName),
			Description:         fmt.Sprintf("Deskripsi Infrastruktur %s", item.SubTypeName),
			InfrastructureCount: item.InfrastructureCount,
		}
	}

	response.Respond(w, http.StatusOK, resp)
}
