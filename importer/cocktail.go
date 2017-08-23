package importer

import (
	"fmt"
	"os"

	"github.com/shrwdflrst/cocktailbot/err"
	"github.com/shrwdflrst/cocktailbot/json"
	"github.com/shrwdflrst/cocktailbot/search"
)

const argrcp = "--recipes"
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
		impr(path)
	} else if name == argsrch {
		terms := args[1:]
		results, e := search.ByIngredient(terms, 0, 10)
		err.Check(e)
		fmt.Println(results)
	}
}

func impr(path string) {
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
	fmt.Println(" --search TERM1 TERM 2...")
}
