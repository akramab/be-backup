package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"mistar-be-go/internal/store"
	"time"

	"github.com/google/uuid"
)

type Infrastructure struct {
	db *sql.DB
}

func NewInfrastructure(db *sql.DB) *Infrastructure {
	return &Infrastructure{db: db}
}

const infrastructureInsert = `INSERT INTO
infrastructure(
	id, sub_type_id, name, details, status, created_at
) values(
	$1, $2, $3, $4, $5, $6	
)
`

func (s *Infrastructure) Insert(ctx context.Context, infrastructure *store.InfrastructureData) error {
	insertStmt, err := s.db.PrepareContext(ctx, infrastructureInsert)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC()
	infrastructureID := uuid.NewString()
	_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
		infrastructureID, infrastructure.SubTypeID, infrastructure.Name,
		infrastructure.Description, infrastructure.Status, createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	infrastructure.ID = infrastructureID

	return nil
}

const infrastructureFindAllQuery = `SELECT i.id, i.name, it.name, it.icon_url, ist.name, ist.icon_url, i.details 
	FROM infrastructure i 
	LEFT JOIN infrastructure_sub_type ist 
	ON i.sub_type_id = ist.id 
	LEFT JOIN infrastructure_type it 
	ON ist.type_id = it.id 
`

func (s *Infrastructure) FindAllInfrastructure(ctx context.Context, name string, typeName string, subTypeName string) ([]*store.InfrastructureData, error) {
	infrastructureList := []*store.InfrastructureData{}
	var queryKeys []string
	var queryParams []interface{}

	query := infrastructureFindAllQuery

	if name != "" {
		queryKeys = append(queryKeys, "InfrastructureName")
		queryParams = append(queryParams, name)
	}

	if typeName != "" {
		queryKeys = append(queryKeys, "TypeName")
		queryParams = append(queryParams, typeName)
	}

	if subTypeName != "" {
		queryKeys = append(queryKeys, "SubTypeName")
		queryParams = append(queryParams, subTypeName)
	}

	for index, key := range queryKeys {
		if index == 0 {
			query = query + "WHERE "
		} else {
			query = query + "AND "
		}

		switch key {
		case "InfrastructureName":
			query = query + fmt.Sprintf(`i.name ILIKE '%%' || $%d || '%%' `, index+1)
		case "TypeName":
			query = query + fmt.Sprintf(`it.name ILIKE '%%' || $%d || '%%' `, index+1)
		case "SubTypeName":
			query = query + fmt.Sprintf(`ist.name ILIKE '%%' || $%d || '%%' `, index+1)
		}
	}

	query = query + `ORDER BY i.created_at DESC `

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		infrastructure := &store.InfrastructureData{}
		err := rows.Scan(
			&infrastructure.ID,
			&infrastructure.Name,
			&infrastructure.TypeName,
			&infrastructure.TypeIconURL,
			&infrastructure.SubTypeName,
			&infrastructure.SubTypeIconURL,
			&infrastructure.Description,
		)
		if err != nil {
			return nil, err
		}
		infrastructureList = append(infrastructureList, infrastructure)
	}

	return infrastructureList, nil
}

func (r *Infrastructure) FindAllInfrastructurePagination(ctx context.Context, offset int, limit int, name string, typeName string, subTypeName string) ([]*store.InfrastructureData, error) {
	infrastructureList := []*store.InfrastructureData{}
	var queryKeys []string
	var queryParams []interface{}

	queryParams = append(queryParams, offset, limit)

	query := infrastructureFindAllQuery

	if name != "" {
		queryKeys = append(queryKeys, "InfrastructureName")
		queryParams = append(queryParams, name)
	}

	if typeName != "" {
		queryKeys = append(queryKeys, "TypeName")
		queryParams = append(queryParams, typeName)
	}

	if subTypeName != "" {
		queryKeys = append(queryKeys, "SubTypeName")
		queryParams = append(queryParams, subTypeName)
	}

	for index, key := range queryKeys {
		if index == 0 {
			query = query + "WHERE "
		} else {
			query = query + "AND "
		}

		switch key {
		case "InfrastructureName":
			query = query + fmt.Sprintf(`i.name ILIKE '%%' || $%d || '%%' `, index+3)
		case "TypeName":
			query = query + fmt.Sprintf(`it.name ILIKE '%%' || $%d || '%%' `, index+3)
		case "SubTypeName":
			query = query + fmt.Sprintf(`ist.name ILIKE '%%' || $%d || '%%' `, index+3)
		}
	}

	query = query + `ORDER BY i.created_at DESC LIMIT $2 OFFSET $1 `

	rows, err := r.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		infrastructure := &store.InfrastructureData{}
		err := rows.Scan(
			&infrastructure.ID,
			&infrastructure.Name,
			&infrastructure.TypeName,
			&infrastructure.TypeIconURL,
			&infrastructure.SubTypeName,
			&infrastructure.SubTypeIconURL,
			&infrastructure.Description,
		)
		if err != nil {
			return nil, err
		}
		infrastructureList = append(infrastructureList, infrastructure)
	}

	return infrastructureList, nil
}
