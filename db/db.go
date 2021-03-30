package db

import (
	"context"
	"github.com/gocql/gocql"
	"os"
	"strings"
	"time"
)

// InitCluster - initialize db connection pool
func InitCluster() *gocql.ClusterConfig {
	cluster := gocql.NewCluster(strings.Join([]string{os.Getenv("CASS_HOST")}, ",")) // potentially you can pass multiple cassandra nodes here
	cluster.Keyspace = os.Getenv("CASS_KEYSPACE")
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("CASS_USERNAME"),
		Password: os.Getenv("CASS_PASSWORD"),
	}
	return cluster
}

// Currency - data point for a single currency at single timestamp
type Currency struct {
	Code    string 		`json:"code"`
	Value   string		`json:"value"`
	PubDate time.Time	`json:"pub_date"`
}

// CurrencyModel - handler for DB session, used to call DB methods
type CurrencyModel struct {
	Session *gocql.Session
}

// InsertAllUnique - insert all unique rows
func (m CurrencyModel) InsertAllUnique(data []Currency) error {
	ctx := context.Background()
	// might be worth doing a batch insert, but the sample size is so small that it doesnt matter
	for _, row := range data {
		err := m.Session.Query("INSERT INTO curfetch.currencies (code, value, pubdate) VALUES (?, ?, ?) IF NOT EXISTS",
			row.Code, row.Value, row.PubDate).WithContext(ctx).Exec()
		if err != nil {
			return err
		}
	}
	ctx.Done()
	return nil
}

// Newest - get newest data points for each currency
func (m CurrencyModel) Newest() (data []Currency, err error){
	var code string
	var value string
	ctx := context.Background()

	// fetch newest date
	var newestDate time.Time
	err = m.Session.Query("SELECT MAX(pubDate) FROM currencies WHERE code = 'USD'").WithContext(ctx).Scan(&newestDate)
	if err != nil {
		return nil, err
	}
	// fetch rows from newestDate
	scanner := m.Session.Query(
		`SELECT code, value FROM currencies WHERE pubDate = ?`,
		newestDate,
		).WithContext(ctx).Iter().Scanner()
	for scanner.Next() {
		err = scanner.Scan(&code, &value)
		if err != nil {
			return nil, err
		}
		data = append(data ,Currency{
			Code:   code,
			Value:   value,
			PubDate: newestDate,
		})

	}
	// close scanner/iterator
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	ctx.Done()
	return data, nil
}

// History - get newest data points for each currency
func (m CurrencyModel) History(code string) (data []Currency, err error){
	var value string
	var pubDate time.Time
	ctx := context.Background()

	// fetch rows for code
	scanner := m.Session.Query(
		`SELECT value, pubDate FROM currencies WHERE code = ?`,
		code,
	).WithContext(ctx).Iter().Scanner()
	for scanner.Next() {
		err = scanner.Scan(&value, &pubDate)
		if err != nil {
			return nil, err
		}
		data = append(data ,Currency{
			Code:   code,
			Value:   value,
			PubDate: pubDate,
		})

	}
	// close scanner/iterator
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	ctx.Done()
	return data, nil
}
