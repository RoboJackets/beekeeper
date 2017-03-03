package backends

import (
	"errors"
	"strconv"
)

// Represents a bin
type DummyBin struct {
	name     string
	capacity uint

	parts map[Component]bool
}

func (b *DummyBin) GetName() string {
	return b.name
}
func (b *DummyBin) GetCapacity() uint {
	return b.capacity
}

func (b *DummyBin) GetParts() []Component {
	comp := make([]Component, 0)
	for v := range b.parts {
		comp = append(comp, v)
	}

	return comp
}

func (b *DummyBin) deletePart(c Component) {
	delete(b.parts, c)
}

// Represents a component.
// Probably will be used in all future backends
type DummyComponent struct {
	id, count uint
	owner     Bin

	name, manufacturer string
}

func (c *DummyComponent) GetName() string {
	return c.name
}
func (c *DummyComponent) GetManufacturer() string {
	return c.manufacturer
}
func (c *DummyComponent) GetId() uint {
	return c.id
}
func (c *DummyComponent) GetBin() Bin {
	return c.owner
}
func (c *DummyComponent) setBin(b Bin) {
	c.owner = b
}

type DummyBackend struct {
	// Map of component ID to component
	// Use this for most lookups
	idLookup map[uint]Component
	// Set of all components (for searches)
	components map[Component]bool
	// List of all the bins
	bins []DummyBin
}

// Makes a very simple backend.
// If specifying a number <= 0, 1 is defaulted to
func NewDummyBackend(numBins uint) Backend {
	if numBins <= 0 {
		numBins = 1
	}

	idLookup := make(map[uint]Component)
	components := make(map[Component]bool)
	newDummy := DummyBackend{
		idLookup:   idLookup,
		components: components,
		bins:       make([]DummyBin, numBins)}

	// Let's make 10 bins
	for i := range newDummy.bins {
		mp := make(map[Component]bool)
		newDummy.bins[i] = DummyBin{
			name:  "A" + strconv.Itoa(i),
			parts: mp,
			// TODO stop hard coding this
			capacity: 3}
	}

	return &newDummy
}

func NewComponent(id, count uint, name, manufacturer string) Component {
	return &DummyComponent{
		id:           id,
		count:        count,
		name:         name,
		manufacturer: manufacturer}
}

// Gets all the components in this dummybackend
func (b *DummyBackend) GetAllComponents() []Component {
	comp := make([]Component, 0)
	for _, bin := range b.bins {
		for _, c := range bin.GetParts() {
			comp = append(comp, c)
		}
	}
	return comp
}

// Adds the component to the bin we think is the most suitable
func (b *DummyBackend) AddComponent(comp Component) (Bin, error) {
	var selectedBin *DummyBin
	for _, v := range b.bins {
		if v.capacity > uint(len(v.parts)) {
			selectedBin = &v
			break
		}
	}
	if selectedBin == nil {
		return nil, errors.New("No more space in bins!")
	}

	// Actually add component to bin
	selectedBin.parts[comp] = true
	comp.setBin(selectedBin)

	// Add lookup pointers for us
	b.idLookup[comp.GetId()] = comp
	b.components[comp] = true
	return selectedBin, nil
}

// Moves a component from it's current bin to a valid one
func (b *DummyBackend) MoveComponent(comp Component, name string) error {
	if comp.GetBin() == nil {
		return errors.New("Comp is not stored in a bin yet!")
	}

	for _, bin := range b.bins {
		if bin.name == name {
			if bin.capacity > uint(len(bin.parts)) {
				comp.GetBin().deletePart(comp)
				// delete(comp.owner.parts, comp)
				comp.setBin(&bin)
				bin.parts[comp] = true
				return nil
			}
			return errors.New("'" + bin.name + "' is over capacity!")
		}
	}
	return errors.New("'" + name + "' was not found!")
}

func (b *DummyBackend) LookupId(id uint) (Component, Bin, error) {
	if component, present := b.idLookup[id]; !present {
		return nil, nil, errors.New("No component found with that ID.")
	} else {
		if component.GetBin() == nil {
			return nil, nil, errors.New("[INTERNAL] The component found has no bin associated with it.")
		}
		return component, component.GetBin(), nil
	}
}
