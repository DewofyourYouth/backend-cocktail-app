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

// API Route handlers
func getListOfCocktails(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var results Cocktails

	for _, v := range cocktails {
		var ingredients []cks.Ingredient
		var directions []cks.Instruction
		fmt.Printf("%d", v.ID)
		db.Where("cocktail_ing_refer = ?", v.ID).Find(&ingredients)
		db.Where("cocktail_dir_refer = ?", v.ID).Find(&directions)
		v.Ingredients = append(v.Ingredients, ingredients...)
		v.Directions = append(v.Directions, directions...)
		results = append(results, v)
	}
	cocktails = results
	fmt.Println(cocktails)
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
	db.Find(&cocktails)
	cs := &cocktails
	fmt.Printf("%T", cs)
	router := mux.NewRouter()
	router.HandleFunc("/", getListOfCocktails).Methods(http.MethodGet)
	router.HandleFunc("/add/", cs.addCocktail).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", router))
}
