package store

import "context"

const (
	InfrastructureStatusNewlyCreated = "NEWLY_CREATED"
	InfrastructureStatusApproved     = "APPROVED"
)

type InfrastructureData struct {
	ID             string
	SubTypeID      int
	SubTypeName    string
	SubTypeIconURL string
	TypeID         int
	TypeName       string
	TypeIconURL    string
	Name           string
	Description    string
	Status         string
}

type Infrastructure interface {
	FindAllInfrastructure(ctx context.Context, name string, typeName string, subTypeName string) ([]*InfrastructureData, error)
	FindAllInfrastructurePagination(ctx context.Context, offset int, limit int, name string, typeName string, subTypeName string) ([]*InfrastructureData, error)
	Insert(ctx context.Context, infrastructure *InfrastructureData) error
}
