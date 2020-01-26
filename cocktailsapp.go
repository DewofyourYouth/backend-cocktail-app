package main

import (
	cks "cocktails/cocktails"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Cocktails is  A slice of Cocktails
type Cocktails []cks.Cocktail
type cocktail cks.Cocktail

var cocktails Cocktails
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*"))
}

func main() {
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	// Migrate the models for DB
	db.AutoMigrate(&cks.Ingredient{})
	db.AutoMigrate(&cks.Instruction{})
	db.AutoMigrate(&cks.Cocktail{})
	// get all cocktail models from DB
	db.Find(&cocktails)
	cs := &cocktails
	// pics := flag.String("pics", "./pics", "/pics")
	flag.Parse()
	var dir string

	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	router := mux.NewRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	router.HandleFunc("/", index)
	router.HandleFunc("/list/", getListOfCocktails).Methods(http.MethodGet)
	router.HandleFunc("/add/", cs.addCocktail).Methods(http.MethodPost)
	router.HandleFunc("/cocktail/{id}/", viewCocktail).Methods(http.MethodGet)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	log.Fatal(srv.ListenAndServe())
}

// Route Handlers for templates
// Main page, displays a list of links to available cocktail recipes
func index(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	db.Close()
	var results cks.Cocktails

	for _, v := range cocktails {
		getIngredientsAndDirections(&v)
		results = append(results, v)
	}
	cocktails = Cocktails(results)
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", cocktails)
}

// displays the detail view of individual recipes
func viewCocktail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var result cks.Cocktail
	db.Where("id = ?", id).First(&result)
	getIngredientsAndDirections(&result)
	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "cocktail.html", result)
}

// gets the ingredients and directions for a given cocktail from the db
func getIngredientsAndDirections(c *cks.Cocktail) {
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var ingredients []cks.Ingredient
	var directions []cks.Instruction
	db.Where("cocktail_ing_refer = ?", c.ID).Find(&ingredients)
	c.Ingredients = append(c.Ingredients, ingredients...)
	db.Where("cocktail_dir_refer = ?", c.ID).Find(&directions)
	c.Directions = append(c.Directions, directions...)
}

// API Route handlers for JSON
func getListOfCocktails(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "db/cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	db.Close()
	var results cks.Cocktails

	for _, v := range cocktails {
		getIngredientsAndDirections(&v)
		results = append(results, v)
	}
	cocktails = Cocktails(results)
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
