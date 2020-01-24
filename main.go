package main

import (
	cks "cocktails/cocktails"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Cocktails is  A slice of Cocktails
type Cocktails []cks.Cocktail
type cocktail cks.Cocktail

var cocktails Cocktails

// Dummy data start
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

// Dummy data end

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

	db.Create(&cks.Cocktail{
		Name:        "White Russian",
		Glass:       "short glass",
		Garnish:     "none",
		Description: "A decadent adult milkshake",
		Ingredients: wrIngredients,
		Directions:  wrDirections})

	db.Create(&cks.Cocktail{
		Name:        "Dry Martini",
		Description: "Preferred beverage of James Bond",
		Glass:       "martini glass",
		Garnish:     "olive or lemon twist",
		Ingredients: dmIngredients,
		Directions:  dmDirections})

	var cktls cks.Cocktails
	var cktl cks.Cocktail
	var ingredients cks.Ingredient
	// var directions cks.Instruction
	db.Model(&cktl).Related(&ingredients)
	type Ing struct {
		Name   string
		Amount float64
		Unit   string
	}
	type Ings []Ing
	type Dir struct {
		Step        int
		Instruction string
	}
	var ings Ings
	var dirs []Dir
	db.Raw("SELECT name, amount, unit FROM ingredients WHERE cocktail_ing_refer=?", 1).Scan(&ings)
	db.Raw("SELECT step, instruction FROM instructions WHERE cocktail_dir_refer=? ORDER BY step", 1).Scan(&dirs)
	// fmt.Println(ings, dirs)
	var ingredientsSlice []cks.Ingredient
	for _, v := range ings {
		fmt.Println(v.Name)
		ingredientsSlice = append(ingredientsSlice, cks.Ingredient{
			Name:             v.Name,
			Amount:           v.Amount,
			Unit:             v.Unit,
			CocktailIngRefer: 1,
		})
	}
	fmt.Println(ingredientsSlice)
	db.Find(&cktls)
	fmt.Println(cktls[0].Name)
	wr := &cktls[0]
	*wr = cks.Cocktail(*wr)
	wr.Ingredients = append(wr.Ingredients, ingredientsSlice...)
	fmt.Println(wr)
	// wr.ingredients = append(wr.Ingredients, ingredientsSlice...)
	fmt.Printf("%v\n", cktls)
	// db.First(&cktl)
	// fmt.Printf("%v\n", cktl)
	// fmt.Println(cks.Cocktails(cktls).MakeCocktailJSON())
	// cks.Cocktails(cktls).Print()
	// cs := &cocktails
	// cocktails = append(cocktails, cktls...)
	router := mux.NewRouter()
	router.HandleFunc("/", getListOfCocktails).Methods(http.MethodGet)
	// router.HandleFunc("/add/", cs.addCocktail).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
