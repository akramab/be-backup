package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mistar-be-go/internal/store"
)

type InfrastructureType struct {
	db *sql.DB
}

func NewInfrastructureType(db *sql.DB) *InfrastructureType {
	return &InfrastructureType{db: db}
}

const infrastructureTypeFindAllQuery = `SELECT ist.id, ist.name, ist.icon_url, it.id, it.name, it.icon_url
	FROM infrastructure_sub_type ist 
	LEFT JOIN infrastructure_type it 
	ON ist.type_id = it.id
	
`

func (s *InfrastructureType) FindAllInfrastructureType(ctx context.Context, typeID int, typeName string, subTypeName string) ([]*store.InfrastructureTypeDataDetail, error) {
	infrastructureTypeList := []*store.InfrastructureTypeDataDetail{}
	var queryKeys []string
	var queryParams []interface{}

	query := infrastructureTypeFindAllQuery

	if typeID != 0 {
		queryKeys = append(queryKeys, "TypeID")
		queryParams = append(queryParams, typeID)
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
		case "TypeID":
			query = query + fmt.Sprintf(`it.id = $%d `, index+3)
		case "TypeName":
			query = query + fmt.Sprintf(`it.name ILIKE '%%' || $%d || '%%' `, index+1)
		case "SubTypeName":
			query = query + fmt.Sprintf(`ist.name ILIKE '%%' || $%d || '%%' `, index+1)
		}
	}

	query = query + `ORDER BY ist.type_id ASC, ist.id ASC`

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		infrastructureType := &store.InfrastructureTypeDataDetail{}
		err := rows.Scan(
			&infrastructureType.SubType.ID,
			&infrastructureType.SubType.SubTypeName,
			&infrastructureType.SubType.IconURL,
			&infrastructureType.ID,
			&infrastructureType.TypeName,
			&infrastructureType.IconURL,
		)
		if err != nil {
			return nil, err
		}
		infrastructureTypeList = append(infrastructureTypeList, infrastructureType)
	}

	return infrastructureTypeList, nil
}

func (r *InfrastructureType) FindAllInfrastructureTypePagination(ctx context.Context, offset int, limit int, typeID int, typeName string, subTypeName string) ([]*store.InfrastructureTypeDataDetail, error) {
	infrastructureTypeList := []*store.InfrastructureTypeDataDetail{}
	var queryKeys []string
	var queryParams []interface{}

	queryParams = append(queryParams, offset, limit)

	query := infrastructureTypeFindAllQuery

	if typeID != 0 {
		queryKeys = append(queryKeys, "TypeID")
		queryParams = append(queryParams, typeID)
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
		case "TypeID":
			query = query + fmt.Sprintf(`it.id = $%d `, index+3)
		case "TypeName":
			query = query + fmt.Sprintf(`it.name ILIKE '%%' || $%d || '%%' `, index+3)
		case "SubTypeName":
			query = query + fmt.Sprintf(`ist.name ILIKE '%%' || $%d || '%%' `, index+3)
		}
	}

	query = query + `ORDER BY ist.type_id ASC, ist.id ASC LIMIT $2 OFFSET $1 `

	log.Println("wowow")
	log.Println(queryParams...)
	rows, err := r.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		infrastructureType := &store.InfrastructureTypeDataDetail{}
		err := rows.Scan(
			&infrastructureType.SubType.ID,
			&infrastructureType.SubType.SubTypeName,
			&infrastructureType.SubType.IconURL,
			&infrastructureType.ID,
			&infrastructureType.TypeName,
			&infrastructureType.IconURL,
		)
		if err != nil {
			return nil, err
		}
		infrastructureTypeList = append(infrastructureTypeList, infrastructureType)
	}

	return infrastructureTypeList, nil
}

const infrastructureSubTypeFindFilterDataQuery = `SELECT ist.id, ist.name, COUNT(i.name)  
	FROM infrastructure i 
	RIGHT JOIN infrastructure_sub_type ist 
	ON i.sub_type_id = ist.id 
`

func (s *InfrastructureType) FindInfrastructureSubTypeFilterData(ctx context.Context, subTypeIdList []int) ([]*store.InfrastructureSubTypeFilterData, error) {
	infrastructureSubTypeFilterDataList := []*store.InfrastructureSubTypeFilterData{}

	query := infrastructureSubTypeFindFilterDataQuery

	if len(subTypeIdList) != 0 {
		queryString := ""
		for idx, id := range subTypeIdList {
			if idx == 0 {
				queryString = queryString + fmt.Sprintf(`%d`, id)
			} else {
				queryString = queryString + fmt.Sprintf(`,%d`, id)
			}
		}
		query = query + fmt.Sprintf(`WHERE ist.id in (%s)`, queryString)
	}

	query = query + `GROUP BY ist.name, ist.id
	ORDER BY ist.id asc`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		infrastructureSubTypeFilterData := &store.InfrastructureSubTypeFilterData{}
		err := rows.Scan(
			&infrastructureSubTypeFilterData.ID,
			&infrastructureSubTypeFilterData.SubTypeName,
			&infrastructureSubTypeFilterData.InfrastructureCount,
		)
		if err != nil {
			return nil, err
		}
		infrastructureSubTypeFilterDataList = append(infrastructureSubTypeFilterDataList, infrastructureSubTypeFilterData)
	}

	return infrastructureSubTypeFilterDataList, nil
}
