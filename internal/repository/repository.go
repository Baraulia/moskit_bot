package repository

import "moskitbot/internal/tv"

type Repository interface {
	GetAll() ([]tv.Line, []string, error)
	GetOne(id int64) (tv.Line, error)
	Create(line tv.Line) (int64, error)
	Delete(id int64) (int64, error)
	Update(id int64, line tv.Line) (int64, error)
}
