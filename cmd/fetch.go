package cmd

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/rzauls/curfetch/db"
	"github.com/spf13/cobra"
	"log"
)

// local command flags
var source string

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch and update currency data",
	Long: `Fetches currency data from RSS feed, upserts data into database`,
	Run: func(cmd *cobra.Command, args []string) {
		fetch()
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringVarP(&source, "source", "s", "https://www.bank.lv/vk/ecb_rss.xml", "rss feed http url")
}

func fetch() {
	feed, err := fetchRssFeed(source)
	if err != nil {
		log.Fatalf("Failed to fetch RSS feed: %v", err)
	}

	if err := db.InitDB("localhost:"); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println(feed.Title)
}

func fetchRssFeed(url string) (feed *gofeed.Feed, err error)  {
	fp := gofeed.NewParser()
	if feed, err := fp.ParseURL(url); err != nil {
		return nil, err
	} else {
		fmt.Println(feed.Title)
		for _ , item := range feed.Items {
			fmt.Println(item.Description)
		}
		return feed, nil
	}
}
