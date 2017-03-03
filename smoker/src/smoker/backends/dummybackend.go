package backends

import (
	"errors"
	"strconv"
)

// Represents a bin
type Bin struct {
	name     string
	capacity uint

	parts map[*Component]bool
}

func (b *Bin) GetName() string {
	return b.name
}
func (b *Bin) GetCapacity() uint {
	return b.capacity
}

func (b *Bin) GetParts() []Component {
	comp := make([]Component, 0)
	for v := range b.parts {
		comp = append(comp, *v)
	}

	return comp
}

// Represents a component
type Component struct {
	id, count uint
	owner     *Bin

	name, manufacturer string
}

func (c *Component) GetName() string {
	return c.name
}
func (c *Component) GetManufacturer() string {
	return c.manufacturer
}
func (c *Component) GetId() uint {
	return c.id
}
func (c *Component) GetBin() *Bin {
	return c.owner
}

type DummyBackend struct {
	// Map of component ID to component
	// Use this for most lookups
	idLookup map[uint]*Component
	// Set of all components (for searches)
	components map[*Component]bool
	// List of all the bins
	bins []Bin
}

// Makes a very simple backend.
// If specifying a number <= 0, 1 is defaulted to
func NewDummyBackend(numBins uint) *DummyBackend {
	if numBins <= 0 {
		numBins = 1
	}

	idLookup := make(map[uint]*Component)
	components := make(map[*Component]bool)
	newDummy := DummyBackend{
		idLookup:   idLookup,
		components: components,
		bins:       make([]Bin, numBins)}

	// Let's make 10 bins
	for i := range newDummy.bins {
		mp := make(map[*Component]bool)
		newDummy.bins[i] = Bin{
			name:  "A" + strconv.Itoa(i),
			parts: mp,
			// TODO stop hard coding this
			capacity: 3}
	}

	return &newDummy
}

func NewComponent(id, count uint, name, manufacturer string) *Component {
	return &Component{
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
func (b *DummyBackend) AddComponent(comp *Component) error {
	var selectedBin *Bin
	for _, v := range b.bins {
		if v.capacity > uint(len(v.parts)) {
			selectedBin = &v
			break
		}
	}
	if selectedBin == nil {
		return errors.New("No more space in bins!")
	}

	// Actually add component to bin
	selectedBin.parts[comp] = true
	comp.owner = selectedBin

	// Add lookup pointers for us
	b.idLookup[comp.id] = comp
	b.components[comp] = true
	return nil
}

// Moves a component from it's current bin to a valid one
func (b *DummyBackend) MoveComponent(comp *Component, name string) error {
	if (comp.owner == nil) {
		return errors.New("Comp is not stored in a bin yet!")
	}

	for _, bin := range b.bins {
		if bin.name == name {
			if bin.capacity < uint(len(bin.parts)) {
				delete(comp.owner.parts, comp)
				comp.owner = &bin
				bin.parts[comp] = true
				return nil
			}
			return errors.New("'" + bin.name + "' is over capacity!")
		}
	}
	return errors.New("'" + name + "' was not found!")
}

func (b *DummyBackend) LookupId(id uint) (*Component, *Bin, error) {
	if component, present := b.idLookup[id]; !present {
		return nil, nil, errors.New("No component found with that ID.")
	} else {
		if component.owner == nil {
			return nil, nil, errors.New("[INTERNAL] The component found has no bin associated with it.")
		}
		return component, component.owner, nil
	}
}
