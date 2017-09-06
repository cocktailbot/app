package err

import "log"

// Check if there was an error and panic
func Check(e error) {
	if e != nil {
		log.Println(e)
		panic(e)
	}
}
