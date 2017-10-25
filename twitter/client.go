package twitter

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Config needed to connect to Twitter
type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

// Create a client
func Create(c Config) *twitter.Client {
	config := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

// Stream tweets containing keywords to callback
func Stream(client *twitter.Client, keywords []string, callback func(interface{})) {
	params := &twitter.StreamFilterParams{
		Track:         keywords,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(params)

	if err != nil {
		panic(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		if !tweet.Retweeted {
			callback(tweet)
		}
	}

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	stream.Stop()
}

// Tweet a message
func Tweet(client *twitter.Client, message string) {
	// Send a Tweet
	// tweet, resp, err := client.Statuses.Update("The time is: "+time.Now().String(), nil)
	_, _, err := client.Statuses.Update(message, nil)
	if err != nil {
		panic(err)
	}
}
