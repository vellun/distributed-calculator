package database

import (
	"context"
	"distributed-calculator/orchestrator/internal/config"
	"distributed-calculator/orchestrator/internal/database"
	"distributed-calculator/orchestrator/pkg/models"
	"errors"
	"fmt"
)

func GetAllExpressions() ([]models.Expression, error) {
	DBParams, err := config.GetDBParams()
	if err != nil {
		return nil, errors.New("Cannont connect to database. Params are wrong")
	}
	conn := database.Connect(DBParams)

	rows, err := conn.Query(context.Background(), "SELECT id, expression, status, started_at, ended_at FROM expressions;")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Query for select expressions from table failed: %v\n", err))
	}

	expressions := []models.Expression{}
	for rows.Next() {
		var exp models.Expression
		err := rows.Scan(&exp.Id, &exp.Expression, &exp.Status, &exp.Started_at, &exp.Ended_at)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error occured while scan expressions: %v\n", err))
		}
		expressions = append(expressions, exp)
	}
	return expressions, nil
}