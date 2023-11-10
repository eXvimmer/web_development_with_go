package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgressConfig struct {
	host, port, user, password, dbname, sslmode string
}

func (p *PostgressConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, p.sslmode)
}

func DefaultPostgresConfig() *PostgressConfig {
	return &PostgressConfig{
		host:     "localhost",
		port:     "5432",
		user:     "goblina",
		password: "jinnythejimbo",
		dbname:   "lenslocked",
		sslmode:  "disable",
	}
}

// OpenDB will open a SQL connection with the provided Postgres configuration.
// Callers of OpenDB need to make sure to close the database connection using
// the `db.close()` method.
func OpenDB(config *PostgressConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("open db, cannot ping: %w", err)
	}
	return db, nil
}
