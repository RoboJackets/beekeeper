
package backends

type Backend interface {
	GetAllComponents() []Component
	MoveComponent(Component, string) error
	LookupId(uint) (Component, Bin, error)
	AddComponent(Component) error
}

type Bin interface {
	GetName() string
	GetCapacity() uint
	GetParts() []Component
	deletePart(Component)
}

type Component interface {
	GetName() string
	GetManufacturer() string
	GetId() uint
	GetBin() Bin
	setBin(Bin)
}
