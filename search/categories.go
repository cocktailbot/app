package search

import "fmt"

// CategoryType that denotes one category
var CategoryType = "category"

// CategoryMapping for index
var CategoryMapping = fmt.Sprintf(`{
	"%s":{
		"properties": {
			"slug" : {
				"type" : "string",
				"index" : "not_analyzed"
			},
			"title": {
				"type":"keyword"
			}
		}
	}
}`, CategoryType)
