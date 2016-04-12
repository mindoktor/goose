package goose

import (
	"database/sql"
	"fmt"
)

// SqlDialect abstracts the details of specific SQL dialects
// for goose's few SQL specific statements
type SqlDialect interface {
	createVersionTableSql(string) string // sql string to create the goose_db_version table
	insertVersionSql(string) string      // sql string to insert the initial version table row
	dbVersionQuery(db *sql.DB, tableName string) (*sql.Rows, error)
}

// drivers that we don't know about can ask for a dialect by name
func dialectByName(d string) SqlDialect {
	switch d {
	case "postgres":
		return &PostgresDialect{}
	}

	return nil
}

////////////////////////////
// Postgres
////////////////////////////

type PostgresDialect struct{}

func (pg PostgresDialect) createVersionTableSql(tableName string) string {
	return fmt.Sprintf(`CREATE TABLE %s (
            	id serial NOT NULL,
                version_id bigint NOT NULL,
                is_applied boolean NOT NULL,
                tstamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );`, tableName)
}

func (pg PostgresDialect) insertVersionSql(tableName string) string {
	return fmt.Sprintf("INSERT INTO %s (version_id, is_applied) VALUES ($1, $2);", tableName)
}

func (pg PostgresDialect) dbVersionQuery(db *sql.DB, tableName string) (*sql.Rows, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT version_id, is_applied from %s ORDER BY id DESC", tableName))

	// XXX: check for postgres specific error indicating the table doesn't exist.
	// for now, assume any error is because the table doesn't exist,
	// in which case we'll try to create it.
	if err != nil {
		return nil, ErrTableDoesNotExist
	}

	return rows, err
}
