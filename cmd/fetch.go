package cmd

import (
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)
// local command flags
var source string

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch and update currency data",
	Long: `Fetches currency data from RSS feed, upserts data into database`,
	Run: func(cmd *cobra.Command, args []string) {

		feed, err := fetchRssFeed(source)
		if err != nil {
			log.Fatalf("Failed to fetch RSS feed: %v", err)
		}
		fmt.Println(feed.Title)
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	// define flags
	fetchCmd.Flags().StringVarP(&source, "source", "s", "https://www.bank.lv/vk/ecb_rss.xml", "rss feed http url")

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
