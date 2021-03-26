package db

import (
	"database/sql"
)

var db *sql.DB

// InitDB - initialize db connection pool
func InitDB(dataSource string) error {
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return err
	}
	return db.Ping()
}

type History struct {
	Code string
	Value float32
	Date string
}

// GetHistory - get historic data of specific currency
func GetHistory(code string) (dataPoints []History, err error){
	// TODO: fetch db rows here

	return []History{
		{Code: code, Value:1.2, Date:"2021-03-27"},
		{Code: code, Value:1.3, Date:"2021-03-26"},
		{Code: code, Value:1.4, Date:"2021-03-25"},
	}, nil
}


// GetNewest - get newest data from database
func GetNewest() (dataPoints []History, err error){
	// TODO: fetch db rows here

	return []History{
		{Code: "USD", Value:1.2, Date:"2021-03-27"},
		{Code: "CAD", Value:1.3, Date:"2021-03-26"},
		{Code: "TET", Value:1.4, Date:"2021-03-25"},
	}, nil
}