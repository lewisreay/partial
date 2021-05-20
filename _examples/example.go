package partial

import (
	"fmt"
	"log"
	"os"
)

type Jedi struct {
	Name            string          `json:"name,omitempty" db:"name"`
	LightSide       bool            `json:"light_side,omitempty" db:"light_side"`
	Padawans        padawans        `json:"padawans,omitempty" db:"padawans"`
	DefeatedEnemies defeatedEnemies `json:"defeated_enemies,omitempty" db:"defeated_enemies"`
	empty           float64         `json:"empty,omitempty" db:"empty"`
}

type (
	defeatedEnemies map[string]bool
	padawans        []string
)

// Implement the Partials interface
func (de defeatedEnemies) Value(i interface{}) (interface{}, error) {
	de, ok := i.(defeatedEnemies)
	if !ok {
		return nil, fmt.Errorf("got %T expected defeatedEnemies", i)
	}
	return de, nil
}

// Implement the Partials interface
func (p padawans) Value(i interface{}) (interface{}, error) {
	p, ok := i.(padawans)
	if !ok {
		return nil, fmt.Errorf("got %T expected padawans", i)
	}
	return p, nil
}

func main() {
	padawans := []string{"Anakin Skywalker", "Ahsoka Tano"}
	jedi := Jedi{
		Name:      "Obi-Wan",
		LightSide: true,
		Padawans:  padawans,
		DefeatedEnemies: map[string]bool{
			"General Grevious": true,
			"Count Dooku":      false,
			"Asajj Ventress":   true,
			"The High Ground":  true,
		},
	}

	v, err := Get(jedi, "db")
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	log.Println(v)
}
