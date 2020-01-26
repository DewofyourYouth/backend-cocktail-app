package cocktails

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// Ingredient is a type that accepts a name: string, amount: float64 and unit: string
type Ingredient struct {
	gorm.Model
	Name             string  `json:"name"`
	Amount           float64 `json:"amount"`
	Unit             string  `json:"unit"`
	CocktailIngRefer uint    `json:"cockail-id"`
}

// Instruction is a type that accepts an int and a string
type Instruction struct {
	gorm.Model
	Step             int    `json:"step"`
	Instruction      string `json:"instruction"`
	CocktailDirRefer uint   `json:"cockail-id"`
}

// Cocktail is a type for a cocktail recipe that accepts Name: string, Ingredients: []Ingredient and Directions: []Instruction
type Cocktail struct {
	gorm.Model
	Name        string        `gorm:"unique;not null" json:"cocktail_name"`
	Description string        `json:"description"`
	Glass       string        `json:"glass"`
	Garnish     string        `json:"garnish"`
	ImageURL    string        `json:"image_url"`
	Ingredients []Ingredient  `gorm:"foreignkey:CocktailIngRefer" json:"ingredients"`
	Directions  []Instruction `gorm:"foreignkey:CocktailDirRefer;" json:"directions"`
}

//Cocktails is
type Cocktails []Cocktail

// Print is a method that prints a formatted Cocktail recipe to the console
func (c Cocktail) Print() {
	fmt.Printf("%s Preparation\n\n", c.Name)
	fmt.Println("Ingredients:")
	for i, v := range c.Ingredients {
		if v.Amount != 1 {
			v.Unit = v.Unit + "s"
		}
		fmt.Printf("\t%d: %v %v %s\n", i+1, v.Amount, v.Unit, v.Name)
	}
	fmt.Println("Instructions:")
	for i, v := range c.Directions {
		fmt.Printf("\t%d: %s\n", i+1, v.Instruction)
	}
	fmt.Printf("\n\n")
}

// Print Cocktails
func (c Cocktails) Print() {
	for i, v := range c {
		fmt.Println("=======", i+1, v.Name, ":", v.Description, "=======")
		fmt.Printf("\n")
		v.Print()
	}
}

// MakeCocktailJSON make cocktails JSON
func (c Cocktail) MakeCocktailJSON() string {
	bs, err := json.Marshal(c)
	if err != nil {
		return `{"message": "invalid cocktail formatting"}`
	}
	return string(bs)
}

// MakeCocktailJSON is
func (c Cocktails) MakeCocktailJSON() string {
	bs, err := json.Marshal(c)
	if err != nil {
		return `{"message": "invalid cocktail formatting"}`
	}
	return string(bs)
}

//UnmarshalCocktailJSON turns cocktail JSON into a string of Cocktail recipes
func UnmarshalCocktailJSON(s string) []Cocktail {
	bs := []byte(s)
	var cocktails []Cocktail
	if err := json.Unmarshal(bs, &cocktails); err != nil {
		fmt.Println(err)
	}
	return cocktails
}

// IntiateModels is a
func IntiateModels() {
	db, err := gorm.Open("sqlite3", "cocktails.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	fmt.Println("connected successfully")
	db.AutoMigrate(Cocktail{})
	db.AutoMigrate(Ingredient{})
	db.AutoMigrate(Instruction{})
}
