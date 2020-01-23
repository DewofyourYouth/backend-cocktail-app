package main

import (
	cks "cocktails/cocktails"
)

// Cocktails is  A slice of Cocktails
type Cocktails []cks.Cocktail

func main() {
	wrIngredients := []cks.Ingredient{
		cks.Ingredient{Name: "vodka", Amount: 2, Unit: "part"},
		cks.Ingredient{Name: "kalua", Amount: 1, Unit: "part"},
		cks.Ingredient{Name: "whole milk", Amount: 2, Unit: "part"},
	}
	wrDirections := []cks.Instruction{
		cks.Instruction{Step: 1, Instruction: "Fill an old fashioned (short) glass with ice"},
		cks.Instruction{Step: 2, Instruction: "Pour over the vodka and kalua"},
		cks.Instruction{Step: 3, Instruction: "Top with whole milk"},
	}
	dmIngredients := []cks.Ingredient{
		cks.Ingredient{Name: "cracked ice", Amount: 2, Unit: "large cubes"},
		cks.Ingredient{Name: "London dry gin", Amount: 2.5, Unit: "ounce"},
		cks.Ingredient{Name: "dry vermouth", Amount: 0.5, Unit: "ounce"},
	}
	dmDirections := []cks.Instruction{
		cks.Instruction{Step: 1, Instruction: "In mixing glass or cocktail shaker filled with ice, combine gin and vermouth."},
		cks.Instruction{Step: 2, Instruction: "Stir well, about 30 seconds, then strain into martini glass."},
		cks.Instruction{Step: 3, Instruction: "Garnish with olive or lemon twist and serve."},
	}

	whiteRussian := cks.Cocktail{Name: "White Russian", Description: "A decadent adult milkshake", Ingredients: wrIngredients, Directions: wrDirections}
	dryMartini := cks.Cocktail{Name: "Dry Martini", Description: "Preferred beverage of James Bond", Ingredients: dmIngredients, Directions: dmDirections}
	// whiteRussian.Print()
	// dryMartini.Print()
	// ctJSON := whiteRussian.MakeCocktailJSON()
	cocktailsSlice := cks.Cocktails{whiteRussian, dryMartini}

	// fmt.Println(cocktailsSlice.MakeCocktailJSON())
	cocktailsSlice.Print()
	// 	fmt.Println(ctJSON)
	// 	fmt.Println(cks.UnmarshalCocktailJSON("[" + ctJSON + "]"))
}
