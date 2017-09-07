package backends

import (
	"errors"
	"strconv"
	"strings"
)

// * Dummy Type Definitions
// ** DummyBin
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

// ** DummyComponent
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
func (c *DummyComponent) GetCount() uint {
	return c.count
}
func (c *DummyComponent) SetCount(u uint) {
	c.count = u
}
func (c *DummyComponent) setBin(b Bin) {
	c.owner = b
}
func (c *DummyComponent) MatchStr(s string) bool {
	totalStr := strings.Join(
		[]string{c.GetName(),
			c.GetManufacturer(),
			strconv.Itoa(int(c.GetId())),
			c.GetBin().GetName()}, " ")
	return strings.Contains(strings.ToLower(totalStr), strings.ToLower(s))
}


// ** DummyCredential
type PasswordCredential struct {
	username, password string
}

func (p *PasswordCredential) GetUsername() string {
	return p.username
}
func (p *PasswordCredential) GetAuth() string {
	return p.password
}

func (p *PasswordCredential) Verify() bool {
	// TODO implement for real
	return p.username == "user" && p.password == "password"
}

func NewCleanDummyCredential() Credential {
	return NewDummyCredential("user", "password")
}

func NewDummyCredential(user, password string) Credential {
	return &PasswordCredential{
		username: user,
		password: password}
}

// ** DummyCredentialManager
// Represents a CredentialManager (storing creds)
type DummyCredentialManager struct {
	creds map[string]Credential
	current Credential
}

func (c *DummyCredentialManager) AddCredential(cred Credential) (error) {
	_, exists := c.creds[cred.GetUsername()]
	if exists {
		return errors.New("Username already exists.");
	}
	c.creds[cred.GetUsername()] = cred;
	return nil;
}
func (c *DummyCredentialManager) RemoveCredential(user string) (error) {
	if len(c.creds) <= 1 {
		return errors.New("Tried to delete last user!");
	}
	delete(c.creds, user);
	return nil;
}
func (c *DummyCredentialManager) DumpUsers() []string {
	keyList := make([]string, len(c.creds))

	i := 0;
	for key, _ := range c.creds {
		keyList[i] = key;
		i++;
	}
	return keyList;
}
func (c *DummyCredentialManager) CurrentUser() (string, error) {
	if c.current == nil {
		return "", errors.New("No current auth");
	}
	return c.current.GetUsername(), nil;
}
func (c *DummyCredentialManager) Login(login Credential) error {
	if candidate, exists := c.creds[login.GetUsername()]; !exists {
		return errors.New("Username not found!")
	} else if candidate.GetUsername() != login.GetUsername() ||
		candidate.GetAuth() != login.GetAuth() {
		// TODO actually verify things
		return errors.New("Wrong auth supplied.")
	} else {
		c.current = candidate;
		return nil;
	}
}
// Make a new credmanager with a default credential and no one logged in
func NewDummyCredentialManager() CredentialManager {
	creds := make(map[string]Credential);
	defaultCreds := NewCleanDummyCredential()
	creds[defaultCreds.GetUsername()] = defaultCreds
	return &DummyCredentialManager{
		creds: creds,
		current: nil}
}

// ** DummyBackend
type DummyBackend struct {
	// Map of component ID to component
	// Use this for most lookups
	idLookup map[uint]Component
	// Set of all components (for searches)
	components map[Component]bool
	// List of all the bins
	bins []DummyBin

	authManager CredentialManager
}

// Makes a very simple backend.
// If specifying a number <= 0, 1 is defaulted to
func NewDummyBackend(numBins uint) Backend {
	if numBins <= 0 {
		numBins = 1
	}

	// TODO remove hard coded number for demo
	numBins = 2

	idLookup := make(map[uint]Component)
	components := make(map[Component]bool)
	newDummy := DummyBackend{
		idLookup:   idLookup,
		components: components,
		bins:       make([]DummyBin, numBins),
		authManager:  NewDummyCredentialManager()}

	// Let's make 10 bins, A00 -> A09
	for i := range newDummy.bins {
		mp := make(map[Component]bool)
		newDummy.bins[i] = DummyBin{
			name:  "A" + "0" + strconv.Itoa(i),
			parts: mp,
			// TODO stop hard coding this
			capacity: 3}
	}

	// hardcoded bins for demo
	// newDummy.bins[0] = DummyBin{
	// 	name:     "C04",
	// 	parts:    make(map[Component]bool),
	// 	capacity: 2}
	// newDummy.bins[1] = DummyBin{
	// 	name:     "B05",
	// 	parts:    make(map[Component]bool),
	// 	capacity: 2}

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
func (b *DummyBackend) GetAllBinNames() []string {
	bins := make([]string, 0)
	for _, bin := range b.bins {
		bins = append(bins, bin.GetName())
	}
	return bins
}
// Adds the component to the bin we think is the most suitable
func (b *DummyBackend) AddComponent(comp Component) (Bin, error) {
	var selectedBin *DummyBin
	for _, v := range b.bins {
		if v.GetCapacity() > uint(len(v.parts)) {
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
	if comp.GetBin().GetName() == name {
		// We are already in the target bin
		return nil
	}

	for _, bin := range b.bins {
		if bin.GetName() == name {
			if bin.GetCapacity() > uint(len(bin.parts)) {
				comp.GetBin().deletePart(comp)
				comp.setBin(&bin)
				bin.parts[comp] = true
				// Don't need to touch b.components or b.idLookup
				return nil
			}
			return errors.New("'" + bin.name + "' is over capacity!")
		}
	}
	return errors.New("Bin '" + name + "' was not found!")
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
func (b *DummyBackend) GeneralSearch(s string) []Component {
	c := make([]Component, 0)
	// Welcome to the worst inventory system on the planet
	for _, comp := range b.GetAllComponents() {
		if comp.MatchStr(s) {
			c = append(c, comp)
		}
	}
	return c
}
func (b *DummyBackend) RemoveComponent(comp Component) error {
	if comp.GetBin() == nil {
		return errors.New("The requested component is not present")
	} else {
		comp.GetBin().deletePart(comp)
		comp.setBin(nil)
		delete(b.idLookup, comp.GetId())
		delete(b.components, comp)
		return nil
	}
}
func (b *DummyBackend) GetCredentialManager() CredentialManager {
	return b.authManager
}
