package plantmodel

import (
	"github.com/skelterjohn/geom"
	"testing"
)

func TestPlantCreate(t *testing.T) {
	species := NewSpecies("Tree")
	plant := NewPlant(species, geom.Coord{9.0, 13.0})

	if plant.age != 0 || plant.radius != 0.0 {
		t.FailNow()
	}
}

func TestPlantEquals(t *testing.T) {
	species := NewSpecies("Tree")
	plant1 := NewPlant(species, geom.Coord{9.0, 13.0})
	plant2 := NewPlant(species, geom.Coord{9.0, 13.0})

	if !plant1.Equals(plant1) {
		t.FailNow()
	}
	if plant1.Equals(plant2) {
		t.FailNow()
	}
}

func TestPlantBounds(t *testing.T) {
	species := NewSpecies("Tree")
	plant := NewPlant(species, geom.Coord{9.0, 13.0})
	bounds := plant.Bounds()
	if bounds.Min.X != 9.0 || bounds.Min.Y != 13.0 {
		t.FailNow()
	}
	if bounds.Max.X != 9.0 || bounds.Max.Y != 13.0 {
		t.FailNow()
	}

	plant.age = 1
	plant.radius = 2.0
	bounds = plant.Bounds()
	if bounds.Min.X != 7.0 || bounds.Min.Y != 11.0 {
		t.FailNow()
	}
	if bounds.Max.X != 11.0 || bounds.Max.Y != 15.0 {
		t.FailNow()
	}
}

func TestPlantGrow(t *testing.T) {
	species := NewSpecies("Tree")
	plant := NewPlant(species, geom.Coord{9.0, 13.0})
	plant.grow()
	if plant.radius != 0.4 || plant.age != 1 || plant.isMature() {
		t.FailNow()
	}
	plant.grow()
	if plant.radius != 0.8 || plant.age != 2 || plant.isMature() {
		t.FailNow()
	}

	// grow for 10 more years
	for i := 0; i < 10; i++ {
		plant.grow()
	}
	if plant.radius != 2.0 || plant.age != 12 || !plant.isMature() {
		t.FailNow()
	}
}
