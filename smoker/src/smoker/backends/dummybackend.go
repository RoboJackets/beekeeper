package backends

// import "fmt"

// Represents a bin
type Bin struct {
	name string
	// Actually a set
	parts *map[Component]bool
}

// Represents a component
type Component struct {
	id    uint
	owner *Bin
}

type DummyBackend struct {
	// Map of component ID to component
	// Use this for most lookups
	idLookup *map[uint]Component
	// Set of all components (for searches)
	components *map[Component]bool
	// List of all the bins
	bins []Bin
}

// Makes a very simple backend.
func NewDummyBackend() *DummyBackend {

	idLookup := make(map[uint]Component)
	components := make(map[Component]bool)
	newDummy := DummyBackend{idLookup: &idLookup, components: &components, bins: make([]Bin, 10)}

	// Let's make 10 bins
	for i := range newDummy.bins {
		mp := make(map[Component]bool)
		newDummy.bins[i] = Bin{
			name:  "Name",
			parts: &mp}
	}

	return &newDummy
}

// TODO remove me
func Add(x, y int) int {
	return (x + y)
}
