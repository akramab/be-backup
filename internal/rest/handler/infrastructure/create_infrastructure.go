package infrastructure

import (
	"encoding/json"
	"log"
	apierror "mistar-be-go/internal/rest/error"
	"mistar-be-go/internal/rest/response"
	"mistar-be-go/internal/store"
	"net/http"
)

type InfrastructureServiceScope struct {
	Capacity    float64 `json:"capacity,omitempty"`
	KKCount     int     `json:"kk_count,omitempty"`
	PeopleCount int     `json:"people_count,omitempty"`
}

type InfrastructureDescription struct {
	Source              string                     `json:"source,omitempty"`
	Ownership           string                     `json:"ownership,omitempty"`
	Stakeholder         string                     `json:"stakeholder,omitempty"`
	WasteType           string                     `json:"waste_type,omitempty"`
	ToiletType          string                     `json:"toilet_type,omitempty"`
	ServiceScope        InfrastructureServiceScope `json:"service_scope,omitempty"`
	ContactPerson       string                     `json:"contact_person,omitempty"`
	MonthlyServiceBill  float64                    `json:"monthly_service_bill,omitempty"`
	WaterProcessingType string                     `json:"water_processing_type,omitempty"`
	WasteProcessingType string                     `json:"waste_processing_type,omitempty"`
}

type CreateInfrastructureRequest struct {
	Name        string                    `json:"name"`
	Type        string                    `json:"type"`
	Description InfrastructureDescription `json:"description"`
}

type CreateInfrastructureResponse struct {
	Message          string `json:"message"`
	InfrastructureID string `json:"infrastructure_id"`
}

func (handler *infrastructureHandler) CreateInfrastructure(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := CreateInfrastructureRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	descriptionJSONString, err := json.Marshal(req.Description)
	if err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	infrastuctureTypeData, err := handler.infrastructureTypeStore.FindAllInfrastructureType(ctx, 0, "", req.Type)
	if err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	if len(infrastuctureTypeData) == 0 {
		response.Error(w, apierror.BadRequestError("type is not valid"))
		return
	}
	subTypeID := infrastuctureTypeData[0].SubType.ID

	newInfrastructureData := &store.InfrastructureData{
		Name:        req.Name,
		SubTypeID:   subTypeID,
		Description: string(descriptionJSONString),
		Status:      store.InfrastructureStatusNewlyCreated,
	}

	if err := handler.infrastructureStore.Insert(ctx, newInfrastructureData); err != nil {
		log.Println("error insert new infrastructure data: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := CreateInfrastructureResponse{
		Message:          "success",
		InfrastructureID: newInfrastructureData.ID,
	}

	response.Respond(w, http.StatusCreated, resp)
}
