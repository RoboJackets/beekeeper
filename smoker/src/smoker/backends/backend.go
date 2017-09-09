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

// Enum definition for CredentialLevel
type CredentialLevel int

const (
	// run `make stringer` when updating this (go is garbage)
	Unknown CredentialLevel = iota
	User
	Admin
)

// GoLang is a flaming pile of trash
const FIRST_CRED CredentialLevel = Unknown
const LAST_CRED CredentialLevel = Admin
const DEFAULT_CRED CredentialLevel = User
const USER_ADMIN CredentialLevel = Admin

// GoLang is worse than writing in raw x86
//go:generate stringer -type=CredentialLevel

type Credential interface {
	// Get username
	GetUsername() string
	// This is password right now, but it really should be an API key or a hashed version.
	GetAuth() string

	GetCredentialLevel() CredentialLevel

	setAuth(string) error
	setCredentialLevel(CredentialLevel) error
}

type CredentialManager interface {
	// Add User
	AddCredential(Credential) error
	// Remove User
	RemoveCredential(Credential) error
	// Updates password for a credential
	UpdateAuth(Credential, string) error
	// Updates permission level for a credential
	UpdatePermission(Credential, CredentialLevel) error
	DumpUsers() ([]Credential, error)
	CurrentUser() (Credential, error)
	Login(Credential) error
}
