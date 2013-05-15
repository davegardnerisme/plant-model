package plantmodel

import (
	"github.com/skelterjohn/geom"
	"github.com/skelterjohn/geom/qtree"
	"log"
	"math/rand"
	"time"
	"errors"
	"os"
	"fmt"
	"bufio"
	"encoding/json"
	"io"
)

const EOL = 0xA

type PlantModel struct {
	qt *qtree.Tree
}

func NewPlantModel(width, height int) *PlantModel {
	return &PlantModel{
		qt: qtree.New(qtree.ConfigDefault(), geom.Rect{geom.Coord{0, 0}, geom.Coord{float64(width), float64(height)}}),
	}
}

// plant quantity random seeds of this species
func (self *PlantModel) RandomSeed(species *Species, quantity int) {
	self.RandomBoundedSeed(species, self.qt.UpperBounds, quantity)
}

// plant quantity random seeds of this species within the specified bounds
func (self *PlantModel) RandomBoundedSeed(species *Species, bounds geom.Rect, quantity int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < quantity; i++ {
		log.Println(self.qt.Bounds.Max)
		x := bounds.Min.X + rand.Float64()*(bounds.Max.X-bounds.Min.X)
		y := bounds.Min.Y + rand.Float64()*(bounds.Max.Y-bounds.Min.Y)
		plant := NewPlant(species, geom.Coord{x, y})
		self.qt.Insert(plant)
		log.Println("Seeding plant", plant)
	}
}

// run simulation
func (self *PlantModel) RunSimulation(years int) {
	for y := 0; y < years; y++ {
		for item := range self.qt.Iterate() {
			plant := item.(*Plant)
			plant.grow()
			if !self.checkDomination(plant) {
				if quantity, bounds := plant.shouldSpawn(); quantity > 0 {
					self.RandomBoundedSeed(plant.species, bounds, quantity)
				}
			}
		}
	}
}

// iterate over all plants
func (self *PlantModel) Iterate() <-chan *Plant {
	ch := make(chan *Plant, self.qt.Count)
	for item := range self.qt.Iterate() {
		ch <- item.(*Plant)
	}
	close(ch)
	return ch
}

// get total number of plants in model
func (self *PlantModel) Size() int {
	return self.qt.Count
}

// iterate all plants in some bounding box (any intersecting this bb)
func (self *PlantModel) IterateBounded(bounds geom.Rect) <-chan *Plant {
	ch := make(chan *Plant)
	go func() {
		intersecting := make(map[qtree.Item]bool, 100)
		self.qt.CollectIntersect(bounds, intersecting)
		for item := range intersecting {
			ch <- item.(*Plant)
		}
		close(ch)
	}()

	return ch
}

// save current state to file
func (self *PlantModel) Save(fn string) error {
	fo, err := os.Create(fn)
    if err != nil {
    	return errors.New(fmt.Sprintf("Failed to open file '%v' (%v)", fn, err))
    }
    // close fo on exit and check for its returned error
    defer func() {
        if err := fo.Close(); err != nil {
            panic(err)
        }
    }()
    w := bufio.NewWriter(fo)

 	for plant := range self.Iterate() {
 		b, err := json.Marshal(plant)
 		if err != nil {
 			return err
 		}
        if _, err := w.Write(b); err != nil {
 			return err
        }
        if err = w.WriteByte(EOL); err != nil {
 			return err
        }
    }

    if err = w.Flush(); err != nil {
    	return err
    }

    log.Println("Saved plant model to", fn)

	return nil
}

// load current state from file
func (self *PlantModel) Load(fn string) error {
    fi, err := os.Open(fn)
    if err != nil {
    	return errors.New(fmt.Sprintf("Failed to open file '%v' (%v)", fn, err))
    }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()	
	r := bufio.NewReader(fi)
	for {
		line, err := r.ReadBytes(EOL)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// reconstruct plant; add to qt
		var plant *Plant
		err = json.Unmarshal(line, &plant)
		if err != nil {
			return err
		}
		log.Println("Loaded", plant)
		self.qt.Insert(plant)
	}

	return nil
}

// check and kill any plants dominated by this plant or this plant if
// any plants dominating it
func (self *PlantModel) checkDomination(plant *Plant) (isDominated bool) {
	isDominated = false
	intersecting := make(map[qtree.Item]bool, 1)
	self.qt.CollectIntersect(plant.Bounds(), intersecting)
	for item := range intersecting {
		nearPlant := item.(*Plant)
		if nearPlant.dominatedBy(plant) {
			log.Println("Removing", nearPlant, "dominated by", plant)
			self.qt.Remove(nearPlant)
		} else if plant.dominatedBy(nearPlant) {
			log.Println("Removing", plant, "dominated by", nearPlant)
			self.qt.Remove(plant)
			isDominated = true
		}
	}

	return
}
