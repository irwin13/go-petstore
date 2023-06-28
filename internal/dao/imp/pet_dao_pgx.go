package imp

import (
	"bytes"

	"github.com/irwin13/go-petstore/internal/dao"
	"github.com/irwin13/go-petstore/pkg/client"
	"github.com/irwin13/go-petstore/pkg/logger"
	"github.com/irwin13/go-petstore/pkg/model/entity"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type PetDaoPgx struct {
	logger   *zap.Logger
	dbClient client.DbClient
}

func NewPetDaoPgx(dbClient client.DbClient) *PetDaoPgx {
	return &PetDaoPgx{
		logger:   logger.GetAppLogger(),
		dbClient: dbClient,
	}
}

func (d *PetDaoPgx) tableName() string {
	return " pet "
}

func (d *PetDaoPgx) Search(filter string) ([]entity.GetPet, error) {

	queryResult := make([]entity.GetPet, 0)
	queryFilter := make([]interface{}, 0)

	var sql bytes.Buffer
	sql.WriteString("SELECT ")
	sql.WriteString("id, ")
	sql.WriteString("name, ")
	sql.WriteString("description ")

	sql.WriteString("FROM ")
	sql.WriteString(d.tableName())
	sql.WriteString("WHERE 1=1 ")

	if filter != "" {
		queryFilter = append(queryFilter, filter)
		sql.WriteString(" AND id = $1")
	}

	d.logger.Debug("Search",
		zap.String("sql", sql.String()),
		zap.Any("parameter", queryFilter),
	)

	conn, err := d.dbClient.GetConnection()
	if err != nil {
		return nil, err
	}

	pgxConn := conn.(*pgx.ConnPool)
	rows, err := pgxConn.Query(sql.String(), queryFilter...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		model := entity.GetPet{}
		if err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.Description,
		); err != nil {
			return nil, err
		}
		queryResult = append(queryResult, model)
	}

	return queryResult, nil
}

func (d *PetDaoPgx) Insert(request entity.InsertPet) (string, error) {

	conn, err := d.dbClient.GetConnection()
	if err != nil {
		return "", err
	}

	pgxConn := conn.(*pgx.ConnPool)

	tx, err := pgxConn.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var randomUUID string

	if request.ID == nil || *request.ID == "" {
		randomUUID = dao.GenerateUuidV4()
	}

	sqlParameters := make([]interface{}, 0)
	var sql bytes.Buffer

	sql.WriteString("INSERT INTO ")
	sql.WriteString(d.tableName())
	sql.WriteString("(")
	sql.WriteString("id, ")
	sql.WriteString("name, ")
	sql.WriteString("description ")
	sql.WriteString(")")
	sql.WriteString(" VALUES ")
	sql.WriteString("(")

	sql.WriteString("$1, ")

	sqlParameters = append(sqlParameters, randomUUID)

	sql.WriteString("$2, ")
	sqlParameters = append(sqlParameters, request.Name)

	sql.WriteString("$3 ")
	sqlParameters = append(sqlParameters, request.Description)

	sql.WriteString(")")

	d.logger.Debug("Insert",
		zap.String("sql", sql.String()),
		zap.Any("parameter", sqlParameters),
	)

	if _, err := tx.Exec(sql.String(),
		sqlParameters...,
	); err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return randomUUID, nil

}

func (d *PetDaoPgx) Update(request entity.UpdatePet) (int64, error) {
	conn, err := d.dbClient.GetConnection()
	if err != nil {
		return -1, err
	}

	pgxConn := conn.(*pgx.ConnPool)

	tx, err := pgxConn.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	var sql bytes.Buffer
	sqlParameters := make([]interface{}, 0)

	sql.WriteString("UPDATE ")
	sql.WriteString(d.tableName())
	sql.WriteString("SET ")

	sql.WriteString("name = $1, ")
	sqlParameters = append(sqlParameters, request.Name)

	sql.WriteString("description = $2 ")
	sqlParameters = append(sqlParameters, request.Description)

	sql.WriteString("WHERE id = $3 ")
	sqlParameters = append(sqlParameters, request.ID)

	d.logger.Debug("Update",
		zap.String("sql", sql.String()),
		zap.Any("parameter", sqlParameters),
	)

	res, err := tx.Exec(sql.String(), sqlParameters...)

	if err != nil {
		return -1, err
	}

	rowCount := res.RowsAffected()

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return rowCount, nil
}

func (d *PetDaoPgx) Delete(request entity.DeletePet) (int64, error) {
	conn, err := d.dbClient.GetConnection()
	if err != nil {
		return -1, err
	}

	pgxConn := conn.(*pgx.ConnPool)

	tx, err := pgxConn.Begin()
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	var sql bytes.Buffer
	sqlParameters := make([]interface{}, 0)

	sql.WriteString("DELETE FROM")
	sql.WriteString(d.tableName())
	sql.WriteString("WHERE id = $1 ")
	sqlParameters = append(sqlParameters, request.ID)

	d.logger.Debug("Delete",
		zap.String("sql", sql.String()),
		zap.Any("parameter", sqlParameters),
	)

	res, err := tx.Exec(sql.String(), sqlParameters...)
	if err != nil {
		return -1, err
	}

	rowCount := res.RowsAffected()

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return rowCount, nil
}
