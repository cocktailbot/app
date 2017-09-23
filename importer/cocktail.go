package main

import (
	"fmt"
	"os"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/json"
	"github.com/cocktailbot/app/search"
)

const argimp = "--import"
const argsrch = "--search"

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]
	name := args[0]

	if name == argimp {
		rpath := args[1]
		cpath := args[2]
		search.CreateIndex(search.Index)
		search.CreateMapping(search.Index, search.RecipeType, search.RecipeMapping)
		search.CreateMapping(search.Index, search.CategoryType, search.CategoryMapping)
		imprt(rpath, search.RecipeType)
		imprt(cpath, search.CategoryType)
	} else if name == argsrch {
		terms := args[1:]
		results, e := search.ByIngredient(terms, 0, 10)
		err.Check(e)
		fmt.Println(results)
	}
}

func imprt(path string, tp string) {
	var items map[string]interface{}
	e := json.Parse(path, &items)
	err.Check(e)
	// fmt.Printf("%#v", items)
	data := items["data"].([]interface{})
	e = search.Save(data, search.Index, tp)
	err.Check(e)
}

func help() {
	fmt.Println("No arguments supplied")
	fmt.Println("\nExamples:")
	fmt.Println(" --recipes PATH_TO_RECIPES_JSON")
	fmt.Println(" --categories PATH_TO_CATEGORIES_JSON")
	fmt.Println(" --search TERM1 TERM 2...")
}
