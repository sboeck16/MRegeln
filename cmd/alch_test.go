package cmd

import (
	"strings"
	"testing"
)

var (
	neededRecipes = map[string]int{
		"MUTAR,FERA": 1,
	}
	disableAlchSearchOnTest = false
)

func TestAlchemicRecipes(t *testing.T) {
	if disableAlchSearchOnTest {
		deb("Recipe - testing disables")
		return
	}
	for k, v := range neededRecipes {
		opts := strings.Split(k, ",")
		if amount := len(searchRecipes(opts)); amount < v {
			t.Error("not enough Recipes for", k,
				"found:", amount, "should be:", v)
		}
	}

}
