package user

import (
	"encoding/json"
	apierror "mistar-be-go/internal/rest/error"
	"mistar-be-go/internal/rest/response"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Name    string `json:"name"`
	Token   string `json:"token"`
	Email   string `json:"email"`
}

func (handler *userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	userData, err := handler.userStore.FindUserByEmail(ctx, req.Email)
	if err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	if userData.Password != req.Password {
		response.Error(w, apierror.ForbiddenError("Incorrect email or password"))
		return
	}

	resp := LoginResponse{
		Message: "Welcome back",
		Name:    cases.Title(language.Und, cases.NoLower).String(userData.Name),
		Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjM4OTliZTBmLWQ5NmItNDliMi05N2EzLTdhOTdlMjQ4NTRlNSIsImVtYWlsIjoidGVzdC1lbWFpbDJAbWFpbC5jb20iLCJpYXQiOjE2ODk3Mzc2NTh9.EHl8V8YrEyGk4idXiCDO2k3Q_rWR-PikfIAUXtIjLyY",
		Email:   userData.Email,
	}

	response.Respond(w, http.StatusCreated, resp)
}
