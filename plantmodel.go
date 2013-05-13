package plantmodel

import (
	"github.com/skelterjohn/geom"
	"github.com/skelterjohn/geom/qtree"
	"log"
	"math/rand"
	"time"
)

type PlantModel struct {
	qt *qtree.Tree
}

func NewPlantModel(width, height int64) *PlantModel {
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

func (self *PlantModel) Iterate() <-chan *Plant {
	ch := make(chan *Plant, self.qt.Count)
	for item := range self.qt.Iterate() {
		ch <- item.(*Plant)
	}
	close(ch)
	return ch
}

func (self *PlantModel) Size() int {
	return self.qt.Count
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
