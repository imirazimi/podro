package adapter

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type SQLDBConfig struct {
	User                  string `koanf:"user"`
	Password              string `koanf:"password"`
	DB                    string `koanf:"db"`
	Host                  string `koanf:"host"`
	Port                  int    `koanf:"port"`
	MaxConnectionLifetime int    `koanf:"max_connection_lifetime"`
	MaxOpenConnections    int    `koanf:"max_open_connections"`
	MaxIdleConnections    int    `koanf:"max_idle_connections"`
}

type SQLDB struct {
	Config SQLDBConfig
	Conn   *sql.DB
}

func NewSQLDB(config SQLDBConfig) *SQLDB {
	return &SQLDB{
		Config: config,
	}
}

func (db *SQLDB) Start() error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.Config.User, db.Config.Password, db.Config.Host, db.Config.Port, db.Config.DB)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	conn.SetConnMaxLifetime(time.Duration(db.Config.MaxConnectionLifetime) * time.Minute)
	conn.SetMaxOpenConns(db.Config.MaxOpenConnections)
	conn.SetMaxIdleConns(db.Config.MaxIdleConnections)

	if err := conn.Ping(); err != nil {
		return err
	}
	db.Conn = conn
	return nil
}

func (db *SQLDB) ShutDown() error {
	return db.Conn.Close()
}
