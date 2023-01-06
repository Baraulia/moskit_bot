package repository

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewMysqlDB(cfg Config) (*sql.DB, error) {
	mysqlCfg := mysql.Config{
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DBName:               cfg.DBName,
		AllowNativePasswords: true,
	}

	path := mysqlCfg.FormatDSN()

	database, err := sql.Open("mysql", path+"&parseTime=true")

	if err != nil {
		return nil, fmt.Errorf("NewMysqlDB:%w", err)
	}

	err = database.Ping()
	if err != nil {
		return nil, fmt.Errorf("NewMysqlDB:%w", err)
	}

	return database, nil
}
