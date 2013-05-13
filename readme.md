# Plant model

Procedural vegetation generator, based on the work in [this paper](http://algorithmicbotany.org/papers/eco.gi2002.html),
as described by [this blog post](http://procworld.blogspot.co.uk/2011/05/forest.html).

Example:

	func main() {
		log.Println("Startup")
		
		model := NewPlantModel(1000, 1000)
		model.RandomSeed(NewSpecies("Pine Tree"), 1)
		
		model.RunSimulation(100)
		
		log.Println("Total plants", model.Size())
		for plant := range model.Iterate() {
			log.Println(plant)
		}
	}

At an early stage thus far.
