package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"moskitbot/internal/tv"
	"moskitbot/pkg/logging"
)

const lineTable = "lines"

type LineRepository struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewLineRepository(database *sql.DB, logger *logging.Logger) *LineRepository {
	return &LineRepository{db: database, logger: logger}
}

func (u *LineRepository) GetOne(id int64) (tv.Line, error) {
	var line tv.Line

	query := fmt.Sprintf("SELECT id, pair, val, description, typ, timeframe FROM %s WHERE id=?", lineTable)

	err := u.db.QueryRow(query, id).Scan(&line.ID, &line.Pair, &line.Val, &line.Description, &line.Typ, &line.Timeframe)

	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Errorf("Line with this id does not exist:%s", err)

			return tv.Line{}, errors.New("line with this id does not exist")
		}

		u.logger.Errorf("Error while scanning for line:%s", err)

		return tv.Line{}, fmt.Errorf("get one line: %w", err)
	}

	return line, nil
}

func (u *LineRepository) GetAll() ([]tv.Line, []string, error) {
	var lines []tv.Line
	var pairs []string

	query := fmt.Sprintf("SELECT id, pair, val, description, typ, timeframe FROM %s", lineTable)

	rows, err := u.db.Query(query)
	if err != nil {
		u.logger.Errorf("Error while executing query for getting list of lines:%s", err)

		return nil, nil, fmt.Errorf("get list of lines: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var line tv.Line
		if err = rows.Scan(&line.ID, &line.Pair, &line.Val, &line.Description, &line.Typ, &line.Timeframe); err != nil {
			u.logger.Errorf("Error while scanning for list of lines:%s", err)

			return nil, nil, fmt.Errorf("get list of lines: %w", err)
		}

		lines = append(lines, line)
	}

	query = fmt.Sprintf("SELECT DISTINCT pair FROM %s", lineTable)

	rows, err = u.db.Query(query)
	if err != nil {
		u.logger.Errorf("Error while executing query for getting list of pairs:%s", err)

		return nil, nil, fmt.Errorf("get list of pairs: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var pair string
		if err = rows.Scan(&pair); err != nil {
			u.logger.Errorf("Error while scanning for list of pairs:%s", err)

			return nil, nil, fmt.Errorf("get list of pairs: %w", err)
		}

		pairs = append(pairs, pair)
	}

	return lines, pairs, nil
}

func (u *LineRepository) Create(line tv.Line) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (pair, val, description, typ, timeframe) VALUES (?, ?, ?, ?, ?)", lineTable)
	result, err := u.db.Exec(query, line.Pair, line.Val, line.Description, line.Typ, line.Timeframe)

	if err != nil {
		u.logger.Errorf("Error while executing query for creating new line:%s", err)

		return 0, fmt.Errorf("create line: %w", err)
	}

	lineID, err := result.LastInsertId()

	if err != nil {
		u.logger.Errorf("Error while getting id of created line:%s", err)

		return 0, fmt.Errorf("get id of created line: %w", err)
	}

	return lineID, nil
}

func (u *LineRepository) Update(id int64, line tv.Line) (int64, error) {
	query := fmt.Sprintf("UPDATE %s SET pair=?, val=?, description+=?, typ=?, timeframe=? WHERE id=?", lineTable)
	result, err := u.db.Exec(query, line.Pair, line.Val, line.Description, line.Typ, line.Timeframe, id)

	if err != nil {
		u.logger.Errorf("Error while executing query for updating line:%s", err)

		return 0, fmt.Errorf("update line: %w", err)
	}

	count, err := result.RowsAffected()

	if err != nil {
		u.logger.Errorf("Error while getting number of affected rows for updating line:%s", err)

		return 0, fmt.Errorf("update line: %w", err)
	}

	return count, nil
}

func (u *LineRepository) Delete(id int64) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", lineTable)
	result, err := u.db.Exec(query, id)

	if err != nil {
		u.logger.Errorf("Error while executing query for deleting line:%s", err)

		return 0, fmt.Errorf("delete line: %w", err)
	}

	count, err := result.RowsAffected()

	if err != nil {
		u.logger.Errorf("Error while getting number of affected rows for deleting line:%s", err)

		return 0, fmt.Errorf("delete line: %w", err)
	}

	return count, nil
}
