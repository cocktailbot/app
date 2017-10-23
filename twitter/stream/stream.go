package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cocktailbot/app/config"
	"github.com/cocktailbot/app/search"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	consumerKey := config.Get("TWITTER_API_KEY")
	consumerSecret := config.Get("TWITTER_SECRET")
	accessToken := config.Get("TWITTER_ACCESS_TOKEN")
	accessSecret := config.Get("TWITTER_ACCESS_TOKEN_SECRET")
	fmt.Println("[" + consumerKey + "]")
	fmt.Println("[" + consumerSecret + "]")
	fmt.Println("[" + accessToken + "]")
	fmt.Println("[" + accessSecret + "]")
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	// Send a Tweet
	// tweet, resp, err := client.Statuses.Update("The time is: "+time.Now().String(), nil)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(tweet)
	// fmt.Println(resp)

	params := &twitter.StreamFilterParams{
		Track:         []string{"kitten"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)

	if err != nil {
		panic(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if !tweet.Retweeted {
			// Convert to json
			byt, err := json.Marshal(tweet)
			if err != nil {
				panic(err)
			}
			// then back into a map
			var dat interface{}
			if err := json.Unmarshal(byt, &dat); err != nil {
				panic(err)
			}

			coll := make([]interface{}, 1, 1)
			coll[0] = dat
			if err := search.Save(coll, search.Index, search.TweetType); err != nil {
				panic(err)
			}
			// fmt.Println(tweet.Text)
		}
	}
	// demux.DM = func(dm *twitter.DirectMessage) {
	// 	fmt.Println(dm.SenderID)
	// }

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	stream.Stop()
}
