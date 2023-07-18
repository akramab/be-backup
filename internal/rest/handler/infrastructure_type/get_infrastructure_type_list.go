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

type GetInfrastructureTypetListRequest struct {
	typeName    string
	subTypeName string
	pageStr     string
	limitStr    string
	page        int
	limit       int
}

func (r *GetInfrastructureTypetListRequest) validate() *apierror.FieldError {
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

type GetInfrastructureTypetListResponse struct {
	ID      int                         `json:"id"`
	Type    string                      `json:"type"`
	IconURL string                      `json:"icon_url"`
	SubType []InfrastructureSubTypeData `json:"sub_type"`
}

type InfrastructureSubTypeData struct {
	ID      int    `json:"id"`
	Name    string `json:"sub_type"`
	IconURL string `json:"icon_url"`
}

func (handler *infrastructureTypeHandler) GetInfrastructureTypeList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := GetInfrastructureTypetListRequest{
		pageStr:     r.URL.Query().Get("page"),
		limitStr:    r.URL.Query().Get("size"),
		typeName:    r.URL.Query().Get("type"),
		subTypeName: r.URL.Query().Get("sub_type"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	offset := req.limit * (req.page - 1)
	infrastructureTypeList, err := handler.infrastructureTypeStore.FindAllInfrastructureTypePagination(ctx, offset, req.limit, 0, req.typeName, req.subTypeName)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := []GetInfrastructureTypetListResponse{}
	infrastructureTypeResponse := GetInfrastructureTypetListResponse{}
	infrastructureSubTypeDataList := []InfrastructureSubTypeData{}
	for _, infrastructureType := range infrastructureTypeList {
		if infrastructureTypeResponse.Type == "" {
			infrastructureTypeResponse = GetInfrastructureTypetListResponse{
				ID:      infrastructureType.ID,
				Type:    infrastructureType.TypeName,
				IconURL: fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, infrastructureType.IconURL),
			}
		}

		if infrastructureTypeResponse.ID == infrastructureType.ID {
			infrastructureSubTypeDataList = append(infrastructureSubTypeDataList, InfrastructureSubTypeData{
				ID:      infrastructureType.SubType.ID,
				Name:    infrastructureType.SubType.SubTypeName,
				IconURL: fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, infrastructureType.SubType.IconURL),
			})
		} else {
			resp = append(resp, GetInfrastructureTypetListResponse{
				ID:      infrastructureTypeResponse.ID,
				Type:    infrastructureTypeResponse.Type,
				IconURL: infrastructureTypeResponse.IconURL,
				SubType: infrastructureSubTypeDataList,
			})

			infrastructureTypeResponse = GetInfrastructureTypetListResponse{
				ID:      infrastructureType.ID,
				Type:    infrastructureType.TypeName,
				IconURL: fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, infrastructureType.IconURL),
			}
			infrastructureSubTypeDataList = []InfrastructureSubTypeData{}
			infrastructureSubTypeDataList = append(infrastructureSubTypeDataList, InfrastructureSubTypeData{
				ID:      infrastructureType.SubType.ID,
				Name:    infrastructureType.SubType.SubTypeName,
				IconURL: fmt.Sprintf("%s/static/%s", handler.apiCfg.Host, infrastructureType.SubType.IconURL),
			})
		}
	}

	resp = append(resp, GetInfrastructureTypetListResponse{
		ID:      infrastructureTypeResponse.ID,
		Type:    infrastructureTypeResponse.Type,
		IconURL: infrastructureTypeResponse.IconURL,
		SubType: infrastructureSubTypeDataList,
	})

	if resp[0].ID == 0 {
		resp = make([]GetInfrastructureTypetListResponse, 0)
	}
	response.Respond(w, http.StatusOK, resp)
}
