package main

import (
	"fmt"
	"os"

	"github.com/cocktailbot/app/err"
	"github.com/cocktailbot/app/json"
	"github.com/cocktailbot/app/search"
)

const argrcp = "--recipes"
const argcat = "--categories"
const argsrch = "--search"

func main() {
	if len(os.Args) < 2 {
		help()
		return
	}

	args := os.Args[1:]
	name := args[0]

	if name == argrcp {
		path := args[1]
		recipes := new(search.Recipes)
		imprt(path, recipes, search.RecipeType)
	} else if name == argcat {
		path := args[1]
		categories := new(search.Categories)
		imprt(path, categories, search.CategoryType)
	} else if name == argsrch {
		terms := args[1:]
		results, e := search.ByIngredient(terms, 0, 10)
		err.Check(e)
		fmt.Println(results)
	}
}

func imprt(path string, items search.Collection, tp string) {
	e := json.Parse(path, &items)
	err.Check(e)

	e = search.Save(items, search.Index, tp)
	err.Check(e)
}

func help() {
	fmt.Println("No arguments supplied")
	fmt.Println("\nExamples:")
	fmt.Println(" --recipes PATH_TO_RECIPES_JSON")
	fmt.Println(" --categories PATH_TO_CATEGORIES_JSON")
	fmt.Println(" --search TERM1 TERM 2...")
}
