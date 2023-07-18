package store

import "context"

type InfrastructureSubTypeData struct {
	ID          int
	SubTypeName string
	IconURL     string
}

type InfrastructureTypeDataDetail struct {
	ID       int
	TypeName string
	IconURL  string
	SubType  InfrastructureSubTypeData
}

type InfrastructureSubTypeFilterData struct {
	ID                  int
	SubTypeName         string
	InfrastructureCount int
}

type InfrastructureType interface {
	FindAllInfrastructureType(ctx context.Context, typeID int, typeName string, subTypeName string) ([]*InfrastructureTypeDataDetail, error)
	FindAllInfrastructureTypePagination(ctx context.Context, offset int, limit int, typeID int, typeName string, subTypeName string) ([]*InfrastructureTypeDataDetail, error)
	FindInfrastructureSubTypeFilterData(ctx context.Context, subTypeIdList []int) ([]*InfrastructureSubTypeFilterData, error)
}
