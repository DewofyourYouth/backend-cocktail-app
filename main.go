package main

import (
	cks "cocktails/cocktails"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Cocktails is  A slice of Cocktails
type Cocktails []cks.Cocktail
type cocktail cks.Cocktail

var cocktails Cocktails

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
	db, dbErr := gorm.Open("sqlite3", "db/cocktails.db")
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
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&cks.Ingredient{})
	db.AutoMigrate(&cks.Instruction{})
	db.AutoMigrate(&cks.Cocktail{})

	db.Create(cks.Cocktail{
		Name:        "Dry Martini",
		Description: "Preferred beverage of James Bond",
		Glass:       "martini glass",
		Garnish:     "olive or lemon twist",
		Ingredients: []cks.Ingredient{
			{Name: "cracked ice", Amount: 2, Unit: "large cubes"},
			{Name: "London dry gin", Amount: 2.5, Unit: "ounce"},
			{Name: "dry vermouth", Amount: 0.5, Unit: "ounce"},
		},
		Directions: []cks.Instruction{
			{Step: 1, Instruction: "In mixing glass or cocktail shaker filled with ice, combine gin and vermouth."},
			{Step: 2, Instruction: "Stir well, about 30 seconds, then strain into martini glass."},
			{Step: 3, Instruction: "Garnish with olive or lemon twist and serve."},
		},
	})

	var cktls cks.Cocktails
	db.Find(&cktls)
	var ingredients cks.Ingredient
	db.Find(&ingredients)
	db.Model(&cktls).Related(&ingredients)
	var dm cks.Cocktail
	db.First(&dm)
	fmt.Println(dm)
	// wr.ingredients = append(wr.Ingredients, ingredientsSlice...)
	fmt.Printf("%v\n", cktls)
	// db.First(&cktl)
	// fmt.Printf("%v\n", cktl)
	// fmt.Println(cks.Cocktails(cktls).MakeCocktailJSON())
	cks.Cocktails(cktls).Print()
	// cs := &cocktails
	cocktails = append(cocktails, cktls...)
	// fmt.Printf("%T", cs)
	// router := mux.NewRouter()
	// router.HandleFunc("/", getListOfCocktails).Methods(http.MethodGet)
	// router.HandleFunc("/add/", cs.addCocktail).Methods(http.MethodPost)
	// log.Fatal(http.ListenAndServe(":8080", router))
}
