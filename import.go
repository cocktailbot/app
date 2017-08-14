package main

import (
	"fmt"
	"os"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/json"
	"github.com/shrwdflrst/cocktailbot/search"
)

const recipes = "--recipes"

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]
	name := args[0]

	if name == recipes {
		path := args[1]
		imprtr(path)
	}
}

func imprtr(path string) {
	var recipes search.Recipes
	e := json.Parse(path, &recipes)
	err.Check(e)
	e = search.Save(recipes)
	err.Check(e)
}

func help() {
	fmt.Println("No arguments supplied")
	fmt.Println("\nExamples:")
	fmt.Println(" --recipes PATH_TO_RECIPES_JSON")
}
