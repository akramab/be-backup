package infrastructuretype

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"mistar-be-go/internal/rest/response"

	apierror "mistar-be-go/internal/rest/error"

	"github.com/go-chi/chi/v5"
)

type GetInfrastructureSubTypeListRequest struct {
	typeName    string
	subTypeName string
	pageStr     string
	limitStr    string
	page        int
	limit       int
}

func (r *GetInfrastructureSubTypeListRequest) validate() *apierror.FieldError {
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

type InfrastructureSubTypeListResponse struct {
	ID      int    `json:"id"`
	SubType string `json:"sub_type"`
	IconURL string `json:"icon_url"`
}

func (handler *infrastructureTypeHandler) GetInfrastructureSubTypeList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GetInfrastructureSubTypeListRequest{
		pageStr:     r.URL.Query().Get("page"),
		limitStr:    r.URL.Query().Get("size"),
		subTypeName: r.URL.Query().Get("sub_type"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	infrastructureSubTypeIDString := chi.URLParam(r, "infrastructure_type_id")
	infrastructureSubTypeID, err := strconv.Atoi(infrastructureSubTypeIDString)
	if err != nil {
		apierror.BadRequestError("infrastructure type id needs to be an integer")
		return
	}

	offset := req.limit * (req.page - 1)
	infrastructureTypeList, err := handler.infrastructureTypeStore.FindAllInfrastructureTypePagination(ctx, offset, req.limit, infrastructureSubTypeID, "", req.subTypeName)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	itemsCount := len(infrastructureTypeList)

	resp := make([]InfrastructureSubTypeListResponse, itemsCount)
	for idx, item := range infrastructureTypeList {
		resp[idx] = InfrastructureSubTypeListResponse{
			ID:      item.SubType.ID,
			SubType: item.SubType.SubTypeName,
			IconURL: fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, item.SubType.IconURL),
		}
	}

	response.Respond(w, http.StatusOK, resp)
}
