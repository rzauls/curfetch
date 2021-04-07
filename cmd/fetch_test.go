package cmd

import (
	"github.com/mmcdole/gofeed"
	"github.com/rzauls/curfetch/db"
	"reflect"
	"testing"
	"time"
)


func Test_parseCurrencyString(t *testing.T) {
	timeNow := time.Now()
	timeFormatted := timeNow.Format("Mon, 02 Jan 2006 03:04:5 -0700")

	validItem := gofeed.Item{
		Description: "AAA 1.1 BBB 22.2 CCC 3.333333",
		Published:   timeFormatted,
		PublishedParsed: &timeNow,
	}

	validItemWant := []db.Currency{
		{
			Code:    "AAA",
			Value:   "1.1",
			PubDate: timeNow,
		},
		{
			Code:    "BBB",
			Value:   "22.2",
			PubDate: timeNow,
		},
		{
			Code:    "CCC",
			Value:   "3.333333",
			PubDate: timeNow,
		},
	}

	tests := []struct {
		name    string
		item    *gofeed.Item
		want    []db.Currency
	}{
		{
			name:    "parses valid string",
			item:    &validItem,
			want:    validItemWant,

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This method errors only when regex compile errors,
			// and since the regex pattern is hardcoded, that should never happen
			got, _ := parseCurrencyString(tt.item)
			for i, point := range got {
				if !reflect.DeepEqual(point, tt.want[i]) {
					t.Errorf("parseCurrencyString() got = %v, want %v", point, tt.want[i])
				}
			}

		})
	}
}

func Test_fetchRssFeed(t *testing.T) {
	feed := gofeed.Feed{}

	type args struct {
		url string
	}
	tests := []struct {
		name     string
		args     args
		wantFeed *gofeed.Feed
		wantErr  bool
	}{
		{
			name:     "successfully fetches remote feed",
			args:     args{url: "http://www.bank.lv/vk/ecb_rss.xml"},
			wantFeed: &feed,
			wantErr:  false,
		},
		{
			name:     "fails to fetch invalid feed",
			args:     args{url: "http://notarealurl.latvia"},
			wantFeed: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFeed, err := fetchRssFeed(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchRssFeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(gotFeed) != reflect.TypeOf(tt.wantFeed) {
				t.Errorf("fetchRssFeed() gotFeed = %v, want %v", gotFeed, tt.wantFeed)
			}
		})
	}
}