package sqldb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Соединяемся с БД
const MAX_RETRY = 5

type DB struct {
	pathToDB string
	dbName   string
	conn     *sql.DB
}

func NewDB(pathToDB, dbName string) *DB {
	return &DB{
		pathToDB: pathToDB,
		dbName:   dbName,
	}
}

func (d *DB) Select(query string, args ...any) (*sql.Rows, error) {
	d.connect()
	defer func() {
		d.close()
	}()

	res, err := d.conn.Query(query, args...)
	if err != nil {
		for i := range MAX_RETRY {
			fmt.Printf("Retry %d of %d", (i + 1), MAX_RETRY)
			d.close()
			res, err = d.conn.Query(query, args...)
			if err == nil {
				return res, nil
			}
		}
		return nil, fmt.Errorf("query (select): %w", err)
	}

	return res, nil
}

func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	d.connect()
	defer func() {
		d.close()
	}()

	res, err := d.conn.Exec(query, args...)
	if err != nil {
		for i := range MAX_RETRY {
			fmt.Printf("Retry %d of %d", (i + 1), MAX_RETRY)
			d.close()
			res, err = d.conn.Exec(query, args...)
			if err == nil {
				return res, nil
			}
		}
		return nil, fmt.Errorf("query (exec): %w", err)
	}

	return res, nil
}

func (d *DB) connect() (bool, error) {
	var err error
	d.conn, err = sql.Open("sqlite3", d.pathToDB+"/"+d.dbName)
	if err != nil {
		fmt.Println("===============")
		fmt.Println("CONNECT", err)
		fmt.Println("===============")

		return false, fmt.Errorf("db connect(open): %w", err)
	}

	if err = d.conn.Ping(); err != nil {
		return false, fmt.Errorf("db connect(ping): %w", err)
	}

	return true, nil
}

func (d *DB) close() error {
	err := d.conn.Close()
	if err != nil {
		return fmt.Errorf("db close: %w", err)
	}

	return nil
}
