package plantmodel

import (
	"fmt"
	"github.com/sdming/gosnow"
	"github.com/skelterjohn/geom"
	"math/rand"
)

var idgen *gosnow.SnowFlake

type Plant struct {
	id       uint64
	position geom.Coord
	age      int
	radius   float64
	species  *Species
}

func NewPlant(species *Species, position geom.Coord) *Plant {
	if idgen == nil {
		idgen, _ = gosnow.Default()
	}
	id, _ := idgen.Next()
	return &Plant{
		id:       id,
		position: position,
		age:      0,
		radius:   0,
		species:  species,
	}
}

func (self *Plant) Bounds() (bounds geom.Rect) {
	return geom.Rect{
		geom.Coord{self.position.X - self.radius, self.position.Y - self.radius},
		geom.Coord{self.position.X + self.radius, self.position.Y + self.radius},
	}
}

func (self *Plant) Equals(oi interface{}) bool {
	if v, ok := oi.(*Plant); ok && v.id == self.id {
		return true
	}
	return false
}

func (self *Plant) String() string {
	return fmt.Sprintf("%v: [%v,%v] %v years old, %v radius", self.id, self.position.X, self.position.Y, self.age, self.radius)
}

// age by 1 year
func (self *Plant) grow() {
	self.age++
	self.radius += self.species.growthRate
	if self.radius > self.species.maxR {
		self.radius = self.species.maxR
	}
}

// check if we're dominated by somet other plant
func (self *Plant) dominatedBy(other *Plant) bool {
	if self.Equals(other) {
		return false
	}
	// bigger plant casts shade @todo could refine
	if self.radius > other.radius {
		return false
	}
	separation := self.position.DistanceFrom(other.position) - self.radius - other.radius
	if separation > 0 {
		return false
	}
	// right now overlap is binary; could factor in how overlapped @todo
	r := rand.Float64()
	var dominated bool
	if self.isMature() {
		dominated = r > self.species.shadeToleranceMature
	} else {
		dominated = r > self.species.shadeToleranceGrowth
	}

	return dominated
}

// are we mature plant?
func (self *Plant) isMature() bool {
	if self.radius >= self.species.maxR {
		return true
	}

	return false
}

// should we spawn?
func (self *Plant) shouldSpawn() (quantity int, bounds geom.Rect) {
	quantity = 0
	bounds = self.species.spawnBounds(self.position)

	if self.isMature() {
		// we may spawn
		if rand.Float64() <= self.species.virility {
			quantity = 1
		}
	}

	return quantity, bounds
}
