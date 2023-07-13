package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	apierror "mistar-be-go/internal/rest/error"
	"mistar-be-go/internal/rest/response"
	"net/http"
	"strconv"
	"strings"
)

type GetInfrastructureListRequest struct {
	name        string
	typeName    string
	subTypeName string
	pageStr     string
	limitStr    string
	page        int
	limit       int
}

func (r *GetInfrastructureListRequest) validate() *apierror.FieldError {
	var err error
	fieldErr := apierror.NewFieldError()

	r.pageStr = strings.TrimSpace(r.pageStr)
	r.limitStr = strings.TrimSpace(r.limitStr)

	if r.pageStr == "" {
		r.pageStr = "1"
	}
	if r.limitStr == "" {
		r.limitStr = "50"
	}

	r.page, err = strconv.Atoi(r.pageStr)
	if err != nil || r.page < 0 {
		fieldErr = fieldErr.WithField("page", "page must be a positive integer")
	}

	r.limit, err = strconv.Atoi(r.limitStr)
	if err != nil || r.limit < 0 {
		fieldErr = fieldErr.WithField("limit", "limit must be a positive integer")
	}

	if len(fieldErr.Fields) != 0 {
		return &fieldErr
	}

	return nil
}

type InfrastructureDataDescription struct {
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

type InfrastructureData struct {
	Name           string                        `json:"name"`
	Type           string                        `json:"type"`
	TypeIconURL    string                        `json:"type_icon_url"`
	SubType        string                        `json:"sub_type"`
	SubTypeIconURL string                        `json:"sub_type_icon_url"`
	Description    InfrastructureDataDescription `json:"description"`
}

type GetAllInfrastructureListResponse struct {
	Page       int                  `json:"page"`
	TotalPages int                  `json:"total_pages"`
	TotalItems int                  `json:"total_items"`
	Items      []InfrastructureData `json:"items"`
}

func (handler *infrastructureHandler) GetAllInfrastructureList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := GetInfrastructureListRequest{
		pageStr:     r.URL.Query().Get("page"),
		limitStr:    r.URL.Query().Get("size"),
		name:        r.URL.Query().Get("name"),
		typeName:    r.URL.Query().Get("type"),
		subTypeName: r.URL.Query().Get("sub_type"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	infrastructureListAll, err := handler.infrastructureStore.FindAllInfrastructure(ctx, req.name, req.typeName, req.subTypeName)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	totalItems := len(infrastructureListAll)

	offset := req.limit * (req.page - 1)
	infrastructureList, err := handler.infrastructureStore.FindAllInfrastructurePagination(ctx, offset, req.limit, req.name, req.typeName, req.subTypeName)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	itemsCount := len(infrastructureList)

	resp := GetAllInfrastructureListResponse{}
	items := make([]InfrastructureData, itemsCount)
	for idx, item := range infrastructureList {
		infrastructureDescription := InfrastructureDataDescription{}
		json.Unmarshal([]byte(item.Description), &infrastructureDescription)

		items[idx] = InfrastructureData{
			Name:           item.Name,
			Type:           item.TypeName,
			TypeIconURL:    fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, item.TypeIconURL),
			SubType:        item.SubTypeName,
			SubTypeIconURL: fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, item.SubTypeIconURL),
			Description:    infrastructureDescription,
		}
	}
	resp.Page = req.page
	resp.TotalItems = totalItems
	resp.TotalPages = int(math.Ceil(float64(totalItems) / float64(req.limit)))
	resp.Items = items

	response.Respond(w, http.StatusOK, resp)
}
