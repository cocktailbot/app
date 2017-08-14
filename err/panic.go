package err

// Check if there was an error and panic
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
