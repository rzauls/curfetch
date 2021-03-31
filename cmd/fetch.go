package cmd

import (
	"github.com/mmcdole/gofeed"
	"github.com/rzauls/curfetch/db"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strings"
)

// local command flags
var source string

// NewFetchCmd represents the fetch command
func NewFetchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fetch",
		Short: "Fetch and update currency data",
		Long:  `Fetches currency data from RSS feed, upserts data into database`,
		Run: func(cmd *cobra.Command, args []string) {
			fetch()
		},
	}
}

// init - initialize command and its flags
func init() {
	fetchCmd := NewFetchCmd()
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&source, "source", "s", "http://www.bank.lv/vk/ecb_rss.xml", "rss feed http url")
}

// fetch - main command action
func fetch() {
	// fetch feed data
	feed, err := fetchRssFeed(source)
	if err != nil {
		log.Fatalf("Failed to fetch RSS feed: %v", err)
	}

	// parse rss feed data
	data, err := parseFeedData(feed)
	if err != nil {
		log.Fatalf("Failed to parse feed data: %v", err)
	}
	// set up db connection
	cluster := db.InitCluster()
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer session.Close()
	currency := db.CurrencyModel{Session: session}

	// post data to db
	err = currency.InsertAllUnique(data)
	if err != nil {
		log.Fatalf("Failed to insert data: %v", err)
	}

	// done
	log.Printf("Fetched and inserted \"%v\"", feed.Title)
}

// fetchRssFeed - fetch rss feed
func fetchRssFeed(url string) (feed *gofeed.Feed, err error) {
	fp := gofeed.NewParser()
	return fp.ParseURL(url)
}

// parseFeedData - unmarshal feed data into structs
func parseFeedData(feed *gofeed.Feed) ([]db.Currency, error) {
	var data []db.Currency
	for _, item := range feed.Items {
		points, err := parseCurrencyString(item)
		if err != nil {
			return nil, err
		}
		data = append(data, points...)
	}
	return data, nil
}

// parseCurrencyString - extract currency data points from feed description
func parseCurrencyString(item *gofeed.Item) ([]db.Currency, error) {
	r, err := regexp.Compile("\\b[A-Z]{3} [0-9]+.[0-9]+\\b")
	if err != nil {
		return nil, err
	}

	var res []string
	var currencies []db.Currency

	for _, result := range r.FindAllString(item.Description, -1) {
		res = strings.Fields(result)
		currencies = append(currencies, db.Currency{
			Code:    res[0],
			Value:   res[1],
			PubDate: *item.PublishedParsed,
		})
	}
	return currencies, nil
}
