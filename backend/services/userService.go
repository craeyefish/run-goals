package services

import (
	"database/sql"
	"log"
	"run-goals/daos"
)

type UserService struct {
	l   *log.Logger
	dao *daos.ActivityDao
}

func NewUserService(l *log.Logger, db *sql.DB) *UserService {
	customerDao := daos.NewActivityDao(l, db)
	return &UserService{
		l:   l,
		dao: customerDao,
	}
}
