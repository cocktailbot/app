package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cocktailbot/app/config"
	"github.com/cocktailbot/app/search"
	"github.com/cocktailbot/app/twitter"
)

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}
	keywords := os.Args[1:]

	fmt.Printf("\nListening for %v\n", keywords)

	config := twitter.Config{
		ConsumerKey:    config.Get("TWITTER_API_KEY"),
		ConsumerSecret: config.Get("TWITTER_SECRET"),
		AccessToken:    config.Get("TWITTER_ACCESS_TOKEN"),
		AccessSecret:   config.Get("TWITTER_ACCESS_TOKEN_SECRET"),
	}
	client := twitter.Create(config)
	twitter.Stream(client, keywords, func(tweet interface{}) {
		// Convert 1 tweet to json
		byt, err := json.Marshal(tweet)
		if err != nil {
			panic(err)
		}
		// then back into a map
		var dat interface{}
		if err := json.Unmarshal(byt, &dat); err != nil {
			panic(err)
		}

		// Save to search index
		coll := make([]interface{}, 1, 1)
		coll[0] = dat
		if err := search.Save(coll, search.Index, search.TweetType); err != nil {
			panic(err)
		}
	})
}

func help() {
	fmt.Println("Access Twitter Streaming API and save tweets containing given keywords")
	fmt.Println("\nExample:")
	fmt.Println("COMMAND keyword1 keyword2")
}
