package json

import (
	"encoding/json"
	"io/ioutil"
)

// Parse a json file
func Parse(path string, v interface{}) error {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return json.Unmarshal(file, &v)
}
