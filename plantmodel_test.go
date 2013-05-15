package plantmodel

import (
	"github.com/skelterjohn/geom"
	"testing"
	"github.com/nu7hatch/gouuid"	
	"fmt"
)

func TestModelCreate(t *testing.T) {
	model := NewPlantModel(100, 100)
	if model.Size() != 0 {
		t.FailNow()
	}
}

func TestModelIterate(t *testing.T) {
	model := NewPlantModel(100, 100)
	species, _ := NewSpecies("plantmodeltest1")
	model.RandomBoundedSeed(species, geom.Rect{geom.Coord{0.0, 0.0}, geom.Coord{100.0, 100.0}}, 100)

	var count int
	for plant := range model.Iterate() {
		plant.grow()
		count++
	}
	if model.Size() != 100 || count != 100 {
		t.FailNow()
	}
}

func TestModelBounds(t *testing.T) {
	model := NewPlantModel(100, 100)
	species, _ := NewSpecies("plantmodeltest2")
	model.RandomBoundedSeed(species, geom.Rect{geom.Coord{10.0, 10.0}, geom.Coord{10.0, 20.0}}, 10)
	model.RandomBoundedSeed(species, geom.Rect{geom.Coord{10.0, 50.0}, geom.Coord{10.0, 60.0}}, 10)
	if model.Size() != 20 {
		t.FailNow()
	}

	var count int
	bounds := geom.Rect{geom.Coord{10.0, 10.0}, geom.Coord{10.0, 20.0}}
	for plant := range model.IterateBounded(bounds) {
		plant.grow()
		count++
	}
	if count != 10 {
		t.FailNow()
	}
}

func TestModelSave(t *testing.T) {
	model := NewPlantModel(100, 100)
	species, _ := NewSpecies("plantmodeltest3")
	model.RandomBoundedSeed(species, geom.Rect{geom.Coord{0.0, 0.0}, geom.Coord{100.0, 100.0}}, 20)

	uuid, _ := uuid.NewV4()
	fn := fmt.Sprintf("/tmp/%v", uuid.String())
	err := model.Save(fn)

	if err != nil {
		t.FailNow()		
	}
}


