package plantmodel

import (
	"github.com/skelterjohn/geom"
	"fmt"
	"errors"
)

var allSpecies map[string]*Species = make(map[string]*Species)

type Species struct {
	// unique textual ID for this species, eg: "pinetree"
	id string
	// how much we grow (radius of the plant) by, at each age in years
	growthRate float64
	// maximum radius we grow to
	maxR float64
	// shade tolerance (0 = none, 1 = total), during growth and mature phases
	shadeToleranceGrowth float64
	shadeToleranceMature float64
	// how likely it is to propagate, at some age (0=won't, 1=will)
	virility float64
	// how far seed likely to travel from plant
	spread float64
}

func NewSpecies(id string) (*Species, error) {
	if _, ok := allSpecies[id]; ok {
		return nil, errors.New(fmt.Sprintf("Species with ID %v already defined", id))
	}
	species := &Species{
		id:                   id,
		growthRate:           0.4,
		maxR:                 2.0,
		shadeToleranceGrowth: 0.2,
		shadeToleranceMature: 0.8,
		virility:             0.4,
		spread:               75.0,
	}
	allSpecies[id] = species

	return species, nil
}

func GetSpecies(id string) (*Species, error) {
	var species *Species
	var ok bool
	if species, ok = allSpecies[id]; !ok {
		return nil, errors.New(fmt.Sprintf("Species with ID %v not defined", id))
	}

	return species, nil
}

func (self *Species) spawnBounds(from geom.Coord) geom.Rect {
	return geom.Rect{
		geom.Coord{from.X - self.spread, from.Y - self.spread},
		geom.Coord{from.X + self.spread, from.Y + self.spread},
	}
}
