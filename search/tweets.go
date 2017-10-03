package search

import "fmt"

// TweetType that denotes one tweet
var TweetType = "tweet"

// TweetMapping for index
var TweetMapping = fmt.Sprintf(`{
	"%s":{
		"properties": {

		}
	}
}`, TweetType)
