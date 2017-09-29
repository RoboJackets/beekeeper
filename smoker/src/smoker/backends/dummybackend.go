package backends

import (
	"encoding/gob"
	"errors"
	"os"
	"strconv"
	"strings"
)

// * Dummy Type Definitions
// ** DummyBin
// Represents a bin
type DummyBin struct {
	Name     string
	Capacity uint

	Parts map[uint]bool
}

func (b *DummyBin) GetName() string {
	return b.Name
}
func (b *DummyBin) GetCapacity() uint {
	return b.Capacity
}

func (b *DummyBin) GetParts() []uint {
	comp := make([]uint, 0)
	for v := range b.Parts {
		comp = append(comp, v)
	}

	return comp
}

func (b *DummyBin) deletePart(c uint) {
	delete(b.Parts, c)
}

// ** DummyComponent
// Represents a component.
// Probably will be used in all future backends
type DummyComponent struct {
	Id, Count uint
	Owner     string

	Name, Manufacturer string
}

func (c *DummyComponent) GetName() string {
	return c.Name
}
func (c *DummyComponent) GetManufacturer() string {
	return c.Manufacturer
}
func (c *DummyComponent) GetId() uint {
	return c.Id
}
func (c *DummyComponent) GetBin() string {
	return c.Owner
}
func (c *DummyComponent) GetCount() uint {
	return c.Count
}
func (c *DummyComponent) SetCount(u uint) {
	c.Count = u
}
func (c *DummyComponent) setBin(b string) {
	c.Owner = b
}
func (c *DummyComponent) MatchStr(s string) bool {
	totalStr := strings.Join(
		[]string{c.GetName(),
			c.GetManufacturer(),
			strconv.Itoa(int(c.GetId())),
			c.GetBin()}, " ")
	return strings.Contains(strings.ToLower(totalStr), strings.ToLower(s))
}

// ** DummyCredential
type PasswordCredential struct {
	Username, Password string
	Level              CredentialLevel
}

func (p *PasswordCredential) GetUsername() string {
	return p.Username
}
func (p *PasswordCredential) GetAuth() string {
	return p.Password
}
func (p *PasswordCredential) GetCredentialLevel() CredentialLevel {
	return p.Level
}
func (p *PasswordCredential) setAuth(auth string) error {
	p.Password = auth
	return nil
}
func (p *PasswordCredential) setCredentialLevel(cred CredentialLevel) error {
	p.Level = cred
	return nil
}

func NewCleanDummyCredential() Credential {
	return NewDummyCredential("user", "password", Admin)
}

func NewDummyCredential(user, password string, level CredentialLevel) Credential {
	if level == Unknown {
		level = DEFAULT_CRED
	}
	return &PasswordCredential{
		Username: user,
		Password: password,
		Level:    level}
}

// ** DummyCredentialManager
// Represents a CredentialManager (storing creds)
type DummyCredentialManager struct {
	Creds   map[string]Credential
	Current Credential
}

var NOPERM_ADMIN string = "Insufficient Permissions. Need: " + USER_ADMIN.String()

func (c *DummyCredentialManager) AddCredential(cred Credential) error {
	if currentUser, err := c.CurrentUser(); err != nil {
		return errors.New("No logged in user.")
	} else if currentUser.GetCredentialLevel() < USER_ADMIN {
		return errors.New(NOPERM_ADMIN)
	}
	_, exists := c.Creds[cred.GetUsername()]
	if exists {
		return errors.New("Username already exists.")
	}
	c.Creds[cred.GetUsername()] = cred
	return nil
}

// Determines the current number of admins (highest privleges)
func (c *DummyCredentialManager) numOfAdministrators() uint {
	var count uint = 0
	for _, cred := range c.Creds {
		if cred.GetCredentialLevel() >= Admin {
			count++
		}
	}
	return count
}
func (c *DummyCredentialManager) RemoveCredential(user Credential) error {
	if currentUser, err := c.CurrentUser(); err != nil {
		return errors.New("No logged in user.")
	} else if currentUser.GetCredentialLevel() < USER_ADMIN {
		return errors.New(NOPERM_ADMIN)
	} else {
		if len(c.Creds) <= 1 {
			return errors.New("Tried to delete last user!")
		}
		if toDel, exists := c.Creds[user.GetUsername()]; !exists {
			return errors.New("No user with name: " + user.GetUsername())
		} else if toDel.GetCredentialLevel() >= Admin && c.numOfAdministrators() <= 1 {
			return errors.New("Attempted to delete last administrator!")
		} else if currentUser.GetUsername() == user.GetUsername() {
			return errors.New("Attempted to delete currently logged in user")
		}
	}

	delete(c.Creds, user.GetUsername())
	return nil
}

func (c *DummyCredentialManager) DumpUsers() ([]Credential, error) {
	if currentUser, err := c.CurrentUser(); err != nil {
		return nil, errors.New("No logged in user.")
	} else if currentUser.GetCredentialLevel() < USER_ADMIN {
		return nil, errors.New(NOPERM_ADMIN)
	}

	valList := make([]Credential, len(c.Creds))
	i := 0
	for _, val := range c.Creds {
		valList[i] = val
		i++
	}
	return valList, nil
}
func (c *DummyCredentialManager) CurrentUser() (Credential, error) {
	if c.Current == nil {
		return nil, errors.New("No current auth")
	}
	return c.Current, nil
}
func (c *DummyCredentialManager) Login(login Credential) error {
	if candidate, exists := c.Creds[login.GetUsername()]; !exists {
		return errors.New("Username not found!")
	} else if candidate.GetUsername() != login.GetUsername() ||
		candidate.GetAuth() != login.GetAuth() {
		// TODO actually verify things
		return errors.New("Wrong auth supplied.")
	} else {
		c.Current = candidate
		return nil
	}
}
func (c *DummyCredentialManager) UpdateAuth(cred Credential, auth string) error {
	if candidate, exists := c.Creds[cred.GetUsername()]; !exists {
		return errors.New("Username not found!")
	} else {
		if currentUser, err := c.CurrentUser(); err == nil {
			if currentUser.GetUsername() != cred.GetUsername() &&
				currentUser.GetCredentialLevel() < USER_ADMIN {
				// If we are not changing our creds and we are not admin, abort.
				return errors.New(NOPERM_ADMIN)
			}
		} else {
			return errors.New("No logged in user.")
		}
		if err := candidate.setAuth(auth); err != nil {
			return err
		}
		return nil
	}
}
func (c *DummyCredentialManager) UpdatePermission(cred Credential, auth CredentialLevel) error {
	if candidate, exists := c.Creds[cred.GetUsername()]; !exists {
		return errors.New("Username not found!")
	} else {
		if currentUser, err := c.CurrentUser(); err == nil {
			if currentUser.GetUsername() != cred.GetUsername() &&
				currentUser.GetCredentialLevel() < USER_ADMIN {
				// If we are not changing our creds and we are not admin, abort.
				return errors.New(NOPERM_ADMIN)
			}
		}
		if candidate.GetCredentialLevel() >= Admin && c.numOfAdministrators() <= 1 {
			return errors.New("Attempted to change permissions of last administrator!")
		}

		if err := candidate.setCredentialLevel(auth); err != nil {
			return err
		}
		return nil
	}
}

// Make a new credmanager with a default credential and no one logged in
func NewDummyCredentialManager() CredentialManager {
	creds := make(map[string]Credential)
	defaultCreds := NewCleanDummyCredential()
	creds[defaultCreds.GetUsername()] = defaultCreds
	return &DummyCredentialManager{
		Creds:   creds,
		Current: nil}
}

// ** DummyBackend
// *** Def
type DummyBackend struct {
	// Map of component ID to component
	// Use this for most lookups
	IdLookup map[uint]Component
	// List of all the bins
	Bins map[string]DummyBin

	AuthManager CredentialManager
}

// *** Constructors
// TODO refactor constructors into their own sections

// Makes a very simple backend.
// If specifying a number <= 0, 1 is defaulted to
func NewDummyBackend(numBins uint) Backend {
	if numBins <= 0 {
		numBins = 1
	}

	idLookup := make(map[uint]Component)
	newDummy := DummyBackend{
		IdLookup:    idLookup,
		Bins:        make(map[string]DummyBin),
		AuthManager: NewDummyCredentialManager()}

	// Let's make 10 bins, A00 -> A09
	for i := 0; i < int(numBins); i++ {
		mp := make(map[uint]bool)
		name := "A" + "0" + strconv.Itoa(i)
		newDummy.Bins[name] = DummyBin{
			Name:  name,
			Parts: mp,
			// TODO stop hard coding this
			Capacity: 3}
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
		Id:           id,
		Count:        count,
		Name:         name,
		Manufacturer: manufacturer}
}

// *** Data Dump Functions

// Gets all the components in this dummybackend
func (b *DummyBackend) GetAllComponents() []Component {
	comp := make([]Component, 0)
	for _, bin := range b.Bins {
		for _, c := range bin.GetParts() {
			comp = append(comp, b.IdLookup[c])
		}
	}
	return comp
}
func (b *DummyBackend) GetAllBinNames() []string {
	bins := make([]string, 0)
	for _, bin := range b.Bins {
		bins = append(bins, bin.GetName())
	}
	return bins
}

// *** Component Modification

// Adds the component to the bin we think is the most suitable
func (b *DummyBackend) AddComponent(comp Component) (Bin, error) {
	var selectedBin *DummyBin
	for _, v := range b.Bins {
		if v.GetCapacity() > uint(len(v.Parts)) {
			selectedBin = &v
			break
		}
	}
	if selectedBin == nil {
		return nil, errors.New("No more space in bins!")
	}

	// Actually add component to bin
	selectedBin.Parts[comp.GetId()] = true
	comp.setBin(selectedBin.GetName())

	// Add lookup pointers for us
	b.IdLookup[comp.GetId()] = comp
	return selectedBin, nil
}

// Moves a component from it's current bin to a valid one
func (b *DummyBackend) MoveComponent(comp Component, name string) error {
	if comp.GetBin() == "" {
		return errors.New("Comp is not stored in a bin yet!")
	}
	if comp.GetBin() == name {
		// We are already in the target bin
		return nil
	}

	for _, bin := range b.Bins {
		if bin.GetName() == name {
			if bin.GetCapacity() > uint(len(bin.Parts)) {
				oldbin := b.Bins[comp.GetBin()]
				oldbin.deletePart(comp.GetId())
				comp.setBin(bin.GetName())
				bin.Parts[comp.GetId()] = true
				// Don't need to touch b.components or b.idLookup
				return nil
			}
			return errors.New("'" + bin.GetName() + "' is over capacity!")
		}
	}
	return errors.New("Bin '" + name + "' was not found!")
}
func (b *DummyBackend) LookupId(id uint) (Component, Bin, error) {
	if component, present := b.IdLookup[id]; !present {
		return nil, nil, errors.New("No component found with that ID.")
	} else {
		if component.GetBin() == "" {
			return nil, nil, errors.New("[INTERNAL] The component found has no bin associated with it.")
		}
		bin := b.Bins[component.GetBin()]
		// TODO take a look at this logic
		return component, &bin, nil
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
	if comp.GetBin() == "" {
		return errors.New("The requested component is not present")
	} else {
		oldBin := b.Bins[comp.GetBin()]
		oldBin.deletePart(comp.GetId())
		comp.setBin("")
		delete(b.IdLookup, comp.GetId())
		return nil
	}
}
func (b *DummyBackend) GetCredentialManager() CredentialManager {
	return b.AuthManager
}

// *** Serialization
func (b *DummyBackend) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(b)
	}
	file.Close()
	return err
}
func (b *DummyBackend) LoadFromFile(path string) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(b)
	}
	file.Close()
	return err
}

// ** Gob Definitions
func RegisterDummyGob() {
	gob.Register(&DummyBackend{})
	gob.Register(&DummyBin{})
	gob.Register(&DummyComponent{})
	gob.Register(&PasswordCredential{})
	gob.Register(&DummyCredentialManager{})
}
