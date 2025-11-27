package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	connectionString string
}

func NewDatabase(connectionString string) *Database {
	return &Database{
		connectionString: connectionString,
	}
}

func (db *Database) GetConnection() (*sql.DB, error) {
	conn, err := sql.Open("postgres", db.connectionString)
	if err != nil {
		log.Println("[Database] Cannot connect to database:", err)
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		log.Println("[Database] Cannot ping database:", err)
		return nil, err
	}

	return conn, nil
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	conn, err := db.GetConnection()
	if err != nil {
		return nil
	}
	defer conn.Close()

	return conn.QueryRow(query, args...)
}

func (db *Database) Query(query string, args ...any) (*sql.Rows, error) {
	conn, err := db.GetConnection()
	if err != nil {
		return nil, err
	}
	// Note: Caller must close rows AND call conn.Close() after use
	// We'll handle this in the service layer

	rows, err := conn.Query(query, args...)
	if err != nil {
		conn.Close()
		log.Println("[Database] Cannot execute query:", err)
		return nil, err
	}

	return rows, nil
}

func (db *Database) Execute(query string, args ...any) (int64, error) {
	conn, err := db.GetConnection()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	result, err := conn.Exec(query, args...)
	if err != nil {
		log.Println("[Database] Cannot execute query:", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[Database] Cannot retrieve affected rows:", err)
		return 0, err
	}

	return rowsAffected, nil
}
