package backends

type Backend interface {
	// Dump all components. Should not be used unless debugging
	// or if you actually want a dump
	GetAllComponents() []Component
	// Moves components to another bin. Error on bin not existing or full
	MoveComponent(Component, string) error
	// Looks up an id. Will return error if the id does not exist
	LookupId(uint) (Component, Bin, error)
	// Adds a component to this backend, will return the suggested bin
	// Error if lack of space
	AddComponent(Component) (Bin, error)
	// General overarching serach of all components.
	// Returns a list of components that match the search
	GeneralSearch(string) []Component
	// Removes a component, and errors if not present.
	RemoveComponent(Component) error
	// Returns a list of all bin names
	GetAllBinNames() []string
	// Return a CredentialManager to handle me Auth
	GetCredentialManager() CredentialManager
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
	GetCount() uint
	SetCount(uint)
	GetBin() Bin
	setBin(Bin)
	// Return true if the input string matches any field in this component
	// Ex: A0 should return true if we are in bin A0, or if we are "resistor A0"
	MatchStr(string) bool
}

type Credential interface {
	// Get username
	GetUsername() string
	// This is password right now, but it really should be an API key or a hashed version.
	GetAuth() string
}

// TODO add permission levels
type CredentialManager interface {
	// Add User
	AddCredential(Credential) error
	// Remove User
	RemoveCredential(Credential) error
	DumpUsers() []string
	CurrentUser() (string, error)
	Login(Credential) error
}
