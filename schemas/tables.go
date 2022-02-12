package schemas

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const Stocks string = `
	CREATE TABLE IF NOT EXISTS  stocks (
    ticker VARCAR(250) NOT NULL PRIMARY KEY,
    description VARCHAR(250),
    etf BOOL NOT NULL,
	expense_ratio REAL,
	insert_ts DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

const Holdings string = `
	CREATE TABLE IF NOT EXISTS holdings (
    etfname VARCAR(250) NOT NULL,
    ticker VARCHAR(250) NOT NULL,
	weight REAL NOT NULL,
	refresh_date DATETIME NOT NULL,
	insert_ts DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (etfname, ticker, refresh_date)
);
`
const Similars string = `
	CREATE TABLE IF NOT EXISTS similars (
    etfname VARCAR(250) NOT NULL,
    ticker VARCHAR(250) NOT NULL,
	refresh_date DATETIME NOT NULL,
	insert_ts DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (etfname, ticker)
);
`

func InitDatabase() error {
	db, err := sql.Open("sqlite3", "etf.db")
	if err != nil {
		return err
	}

	//close the database connection
	defer db.Close()
	// create all the tables
	tables := []string{Stocks, Holdings, Similars}
	for _, val := range tables {
		err := InitTables(db, val)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func InitTables(db *sql.DB, query string) error {
	// log.Println("table being created: ", query)
	_, err := db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

func NewDbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "etf.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
