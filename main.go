package main

import (
	cks "cocktails/cocktails"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Cocktails is  A slice of Cocktails
type Cocktails []cks.Cocktail
type cocktail cks.Cocktail

var cocktails Cocktails

var wrIngredients []cks.Ingredient = []cks.Ingredient{
	cks.Ingredient{Name: "vodka", Amount: 2, Unit: "part"},
	cks.Ingredient{Name: "kalua", Amount: 1, Unit: "part"},
	cks.Ingredient{Name: "whole milk", Amount: 2, Unit: "part"},
}

var wrDirections []cks.Instruction = []cks.Instruction{
	cks.Instruction{Step: 1, Instruction: "Fill an old fashioned (short) glass with ice"},
	cks.Instruction{Step: 2, Instruction: "Pour over the vodka and kalua"},
	cks.Instruction{Step: 3, Instruction: "Top with whole milk"},
}

var dmIngredients []cks.Ingredient = []cks.Ingredient{
	cks.Ingredient{Name: "cracked ice", Amount: 2, Unit: "large cubes"},
	cks.Ingredient{Name: "London dry gin", Amount: 2.5, Unit: "ounce"},
	cks.Ingredient{Name: "dry vermouth", Amount: 0.5, Unit: "ounce"},
}

var dmDirections []cks.Instruction = []cks.Instruction{
	cks.Instruction{Step: 1, Instruction: "In mixing glass or cocktail shaker filled with ice, combine gin and vermouth."},
	cks.Instruction{Step: 2, Instruction: "Stir well, about 30 seconds, then strain into martini glass."},
	cks.Instruction{Step: 3, Instruction: "Garnish with olive or lemon twist and serve."},
}

// API Route handlers
func getListOfCocktails(w http.ResponseWriter, r *http.Request) {
	cktls := cks.Cocktails(cocktails).MakeCocktailJSON()
	w.Header().Set("Content-Type", "application-json")
	fmt.Fprintf(w, cktls)
}

func (cs *Cocktails) addCocktail(w http.ResponseWriter, r *http.Request) {
	var c cks.Cocktail
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*cs = append(*cs, c)
	newCocktail := &c
	db, dbErr := gorm.Open("sqlite3", "cocktails.db")
	if dbErr != nil {
		panic("failed to connect")
	}
	defer db.Close()
	db.Create(newCocktail)
	bs, err1 := json.Marshal(newCocktail)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", bs)
}

func main() {
	db, err := gorm.Open("sqlite3", "cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&cks.Ingredient{})
	db.AutoMigrate(&cks.Instruction{})
	db.AutoMigrate(&cks.Cocktail{})

	whiteRussian := cks.Cocktail{Name: "White Russian", Description: "A decadent adult milkshake", Ingredients: wrIngredients, Directions: wrDirections}
	dryMartini := cks.Cocktail{Name: "Dry Martini", Description: "Preferred beverage of James Bond", Ingredients: dmIngredients, Directions: dmDirections}
	// db.Create(&whiteRussian)
	whiteRussian.Print()
	dryMartini.Print()
	cocktailsSlice := Cocktails{whiteRussian, dryMartini}
	for _, v := range cocktailsSlice {
		db.Create(&v)
	}
	fmt.Println(cks.Cocktails(cocktailsSlice).MakeCocktailJSON())
	cks.Cocktails(cocktailsSlice).Print()
	cs := &cocktails
	cocktails = append(cocktails, cocktailsSlice...)
	router := mux.NewRouter()
	router.HandleFunc("/", getListOfCocktails).Methods(http.MethodGet)
	router.HandleFunc("/add/", cs.addCocktail).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))

}
