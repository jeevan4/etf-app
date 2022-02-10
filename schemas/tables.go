package schemas

const Stocks string = `
	CREATE TABLE IF NOT EXISTS  stocks (
    ticker VARCAR(250) NOT NULL PRIMARY KEY,
    description VARCHAR(250),
    etf BOOL NOT NULL,
	insert_ts DATETIME DEFAULT CURRENT_TIMESTAMP,	
);
`

const Holdings string = `
	CREATE TABLE IF NOT EXISTS holdings (
    etfname VARCAR(250) NOT NULL,
    ticker VARCHAR(250) NOT NULL,
    weight REAL NOT NULL,
	refresh_date DATETIME NOT NULL,
	insert_ts DATETIME DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (etfname, ticker)
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
