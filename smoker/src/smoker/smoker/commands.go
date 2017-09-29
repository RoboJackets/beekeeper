package main

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/getwe/figlet4go"
	"os"
	"smoker/backends"
	"strconv"
	"strings"
	"text/tabwriter"
)

// * Code
// ** Variables
// Colors for different prompt
var scanColor = color.New(color.FgYellow).SprintFunc()
var countColor = color.New(color.FgGreen).SprintFunc()
var binColor = color.New(color.FgBlue).SprintFunc()
var authColor = color.New(color.FgRed).Add(color.Bold).SprintFunc()
var moveColor = binColor

// ** Command Definitions

func initCommands() {
	backends.RegisterDummyGob()

	commands = make(map[string]func([]string, backends.Backend))
	commands["help"] = replHelp
	commands["h"] = replHelp

	commands["dump"] = replDump
	commands["d"] = replDump

	commands["bins"] = replBins
	commands["b"] = replBins

	commands["grep"] = replGrep
	commands["g"] = replGrep

	commands["scan"] = replScan
	commands["s"] = replScan

	commands["save"] = replSave
	commands["load"] = replLoad

	commands["rm"] = replRm
	commands["r"] = replRm

	commands["mv"] = replMove
	commands["m"] = replMove

	commands["u"] = replUpdate
	commands["update"] = replUpdate

	commands["moo"] = replMoo

	commands["login"] = replLogin
	commands["useradd"] = replAddUser
	commands["adduser"] = replAddUser
	commands["userdel"] = replDeleteUser
	commands["deluser"] = replDeleteUser
	commands["who"] = replListUsers
	commands["whoami"] = replWhoAmI
	commands["passwd"] = replChangeAuth
	commands["chperm"] = replChangePermissionLevel

	commands["w"] = replWelcome
	commands["welcome"] = replWelcome

	commands["quit"] = replQuit
	commands["bye"] = replQuit
	commands["q"] = replQuit
}

// ** High Level Funcs
// Runs a repl command by keying into the command map
func runCommand(prompt string, backend backends.Backend) {
	cmds := strings.Fields(prompt)

	if len(cmds) == 0 {
		return
	}

	if v, p := commands[cmds[0]]; !p {
		fmt.Println("'" + cmds[0] + "' is not a valid command.")
	} else {
		v(cmds[1:], backend)
	}
}

func replHelp(args []string, b backends.Backend) {

	// UserAdmin Help page
	if len(args) >= 1 && args[0] == "useradmin" {
		fmt.Println(`This manual page covers user managment commands.

These include deletion/creation of users and setting user permissions.
Most of these commands require admin access, but some (eg: whoami) do not.

List of Commands:`)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "login <USER>\tLogs into the specified user.")
		fmt.Fprintln(w, "useradd <USER>\tAdds the specified user to the system.")
		fmt.Fprintln(w, "userdel <USER>\tRemoves the specified user to the system.")
		fmt.Fprintln(w, "who\tLists all users and permissions on the system.")
		fmt.Fprintln(w, "whoami\tLists the current user and permission level.")
		fmt.Fprintln(w, "passwd (<USER>)\tChanges the password for a user. With no args, this is the current user.")
		fmt.Fprintln(w, "chperm <USER>\tChanges permissions for a user.")
		w.Flush()
		return
	} else {
		// Generic Help Page
		fmt.Println(`Welcome to the Smoker Help Page.

Commands with a (*) have a no-args scanning mode.

List of Commands:`)
		// Observe how the b's and the d's, despite appearing in the
		// second cell of each line, belong to different columns.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "(h)elp\tPrints this help message.")
		fmt.Fprintln(w, "(w)elcome\tPrints a overview message, helpful to beginners.")
		fmt.Fprintln(w, "(q)uit\tQuit's the current repl mode. If at top level, quit.")
		fmt.Fprintln(w, "(d)ump\tDumps information for all components.")
		fmt.Fprintln(w, "(r)m* <ID> [<ID>, ...]\tDeletes one or more components, by ID.")
		fmt.Fprintln(w, "(m)v* <ID> <BIN>\tMoves <ID> to <BIN> if possible.")
		fmt.Fprintln(w, "(u)pdate* <ID> <COUNT>\tUpdates the count of ID to COUNT.")
		fmt.Fprintln(w, "(b)ins\tPrints a list of all bins available.")
		fmt.Fprintln(w, "(g)rep <search>\tGreps all information in every component.")
		fmt.Fprintln(w, "(s)can*\tLaunches the interactive scanner interface to add/identify parts.\n\t  Takes in Component IDs, which can be printed with a scanner\n\t  (q)uit to exit scanning mode.")
		fmt.Fprintln(w, "save <file>\tSaves the current inventory database to a file")
		fmt.Fprintln(w, "load <file>\tLoads the current inventory database from a file")
		w.Flush()
		fmt.Println("")
		fmt.Println("Additional pages:")
		fmt.Println("\t\t\thelp useradmin")
	}

}

func printDump(c []backends.Component) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID"+"\t"+"BIN"+"\t"+"NAME"+"\t"+"MANUFACTURER"+"\t"+"COUNT"+"\t")
	for _, v := range c {
		fmt.Fprintln(w, strconv.Itoa(int(v.GetId()))+"\t"+v.GetBin()+"\t"+v.GetName()+"\t"+v.GetManufacturer()+"\t"+strconv.Itoa(int(v.GetCount()))+"\t")
	}
	w.Flush()
}

// Encode via Gob to file
func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Decode Gob file
func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

// ** Command Definitions
// *** Item Commands
type UserG struct {
	Name, Pass string
	idLookup   map[uint]uint
}

func replDump(s []string, b backends.Backend) {
	c := b.GetAllComponents()

	if len(c) == 0 {
		fmt.Println("No data is present.")
		return
	}
	printDump(c)
}

func replSave(s []string, b backends.Backend) {
	// TODO warn before overwriting
	if len(s) < 1 {
		fmt.Println("Please provide a file to save to.")
	} else {
		err := b.SaveToFile(s[0])
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
	}
}
func replLoad(s []string, b backends.Backend) {
	if len(s) < 1 {
		fmt.Println("Please provide a file to load from.")
	} else {
		err := b.LoadFromFile(s[0])
		if err != nil {
			fmt.Println("Error: " + err.Error())
		}
	}
}

func replBins(s []string, b backends.Backend) {
	allBins := b.GetAllBinNames()
	for i, _ := range allBins {
		fmt.Println(allBins[i])
	}
}

// TODO refactor replScan to use shared read functions
func replScan(s []string, b backends.Backend) {
	for {
		idInt, err := readUint(scanColor("scan*> "), IDWarning)
		if err != nil {
			break
		}

		componentId := idInt
		component, bin, err := b.LookupId(componentId)
		if err != nil {
			// Add a new item
			comp, err := genComponent(componentId)
			if err != nil {
				fmt.Println("Abort adding component.")
				continue
			}
			bin, err = b.AddComponent(comp)
			if err != nil {
				fmt.Println("Failed adding part: " + err.Error())
				continue
			}

			// Display added part
			printGenericInfo(comp, bin)

			for {
				newBin, err := readRaw(moveColor("move?> "))
				if err != nil || len(newBin) == 0 || newBin == "quit" || newBin == "q" {
					if err != nil {
						fmt.Println()
					}
					// we got nothing to move to...
					break
				}
				if err := b.MoveComponent(comp, newBin); err == nil {
					// Success!
					break
				} else {
					fmt.Println("Could not move component: " + err.Error())
					continue
				}
			}
		} else {
			// Display found item
			printGenericInfo(component, bin)
		}
	}
}

func printGenericInfo(c backends.Component, b backends.Bin) {
	w := color.New(color.FgWhite).Add(color.Bold)
	r := color.New(color.FgGreen).Add(color.Bold)
	printBin(b.GetName())
	r.Print("Count:\t")
	w.Println(strconv.Itoa(int(c.GetCount())))
}

func printBin(s string) {
	w := color.New(color.FgWhite).Add(color.Bold)
	r := color.New(color.FgGreen).Add(color.Bold)
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render(s)
	w.Println(renderStr)
	r.Print("Bin:\t")
	w.Println(s)
}

func replMoo(s []string, b backends.Backend) {
	moo, _ := base64.StdEncoding.DecodeString("IF9fX19fX19fIA0KPCBNb28uLi4" +
		"gPg0KIC0tLS0tLS0tIA0KICAgXA0KICAgIFwgICAgICAgICAgICAgIC4uLi4gICAgICAgDQ" +
		"ogICAgICAgICAgIC4uLi4uLi4uICAgIC4gICAgICANCiAgICAgICAgICAuICAgICAgICAgI" +
		"CAgLiAgICAgIA0KICAgICAgICAgLiAgICAgICAgICAgICAuICAgICAgDQouLi4uLi4uLi4g" +
		"ICAgICAgICAgICAgIC4uLi4uLi4NCi4uLi4uLi4uLi4uLi4uLi4uLi4uLi4uLi4uLi4uLg0K")
	fmt.Println(string(moo))
}

func replGrep(s []string, b backends.Backend) {
	args := strings.Join(s, " ")
	if len(args) == 0 {
		fmt.Println("Please enter a search query to grep for.")
		return
	}
	components := b.GeneralSearch(args)
	printDump(components)
}

func replMove(args []string, b backends.Backend) {
	if len(args) == 0 {
		// Interactive mode
		for {
			if id, err := readUint(scanColor("move ID> "), IDWarning); err != nil {
				// User quit, or error reading
				return
			} else if component, _, err := b.LookupId(id); err != nil {
				fmt.Println("'" + strconv.Itoa(int(id)) + "' " + IDWarning)
			} else {
			READBIN:
				if bin, err := readStringLoop(binColor("bin> ")); err != nil {
					return
				} else if err := b.MoveComponent(component, bin); err != nil {
					fmt.Println("Error: " + err.Error())
					// Read the bin again!
					goto READBIN
				}
			}
			// (else) success!
		}
	} else {
		if len(args) != 2 {
			fmt.Println("mv needs 2 or 0 arguments.")
			return
		}
		if id, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("'" + args[0] + "' is not a valid ID!")
			return
		} else if component, _, err := b.LookupId(uint(id)); err != nil {
			fmt.Println("No component with id '" + strconv.Itoa(id) + "' was found.")
			return
		} else if err := b.MoveComponent(component, args[1]); err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
		// We completed successfully!
	}
}

const ErrorDeleteID string = "[INTERNAL] An internal error occurred when deleting an element. Partial deletion probably occured."

func replRm(args []string, b backends.Backend) {
	if len(args) == 0 {
		// Interactive mode
		for {
			if id, err := readUint(scanColor("rm*> "), IDWarning); err != nil {
				// User quit, or error reading
				return
			} else if component, _, err := b.LookupId(id); err != nil {
				fmt.Println("'" + strconv.Itoa(int(id)) + "' " + IDWarning)
			} else if err := b.RemoveComponent(component); err != nil {
				fmt.Println(ErrorDeleteID)
				return
			}
			// (else) Success!
		}
	} else {
		components := make([]backends.Component, 0)
		for _, v := range args {
			if id, err := strconv.Atoi(v); err != nil {
				fmt.Println("'" + v + "' " + IDWarning)
				return
			} else if comp, _, err := b.LookupId(uint(id)); err != nil {
				fmt.Println("'" + v + "' " + IDWarning)
				return
			} else {
				components = append(components, comp)
			}
		}

		for _, comp := range components {
			if err := b.RemoveComponent(comp); err != nil {
				fmt.Println(ErrorDeleteID)
				return
			}
		}
	}
}

// Function for updating the count of a component
func replUpdate(args []string, b backends.Backend) {
	if len(args) == 0 {
		// Interactive mode
		for {
			if id, err := readUint(scanColor("update*> "), IDWarning); err != nil {
				// User quit, or error reading
				return
			} else if component, _, err := b.LookupId(id); err != nil {
				fmt.Println("'" + strconv.Itoa(int(id)) + "' " + IDWarning)
			} else if count, err := readUint(countColor("count> "), CountWarning); err != nil {
				return
			} else {
				// TODO consider relative counts
				component.SetCount(uint(count))
			}
		}
	} else {
		if len(args) != 2 {
			fmt.Println("Update must take 2 or no arguments.")
			return
		}

		if i, err := strconv.ParseUint(args[0], 10, 32); err != nil {
			fmt.Println("'" + args[0] + "' is not a valid ID.")
		} else if component, _, err := b.LookupId(uint(i)); err != nil {
			fmt.Println("'" + args[0] + "' is not a present component.")
		} else if j, err := strconv.ParseUint(args[1], 10, 32); err != nil {
			fmt.Println("'" + args[1] + "' is not a valid count.")
		} else {
			component.SetCount(uint(j))
		}
	}
}

func replWelcome(args []string, b backends.Backend) {
	fmt.Println("Welcome to Smoker!")
	fmt.Println()
	fmt.Println("Smoker is the CLI frontend to the BeeKeeper inventory suite.")
	fmt.Println("It offers quick and easy access to all functionality beekeeper provides, while staying out of your way as much as possible.")
	fmt.Println()
	fmt.Println("A good place to get started is the 'scan' command, which opens a seperate scanning REPL, allowing you to add new items if available, and print locations of items if found.")
	fmt.Println()
	fmt.Println("For documentation, type 'help'")
}

// Queries the user for the info required to make a component
func genComponent(id uint) (backends.Component, error) {
	// Need Count, Name, and Manufacturer
	if name, err := readStringLoopRaw("Part Name> ", false); err != nil {
		return nil, errors.New("Error Reading Component.")
	} else if man, err := readStringLoopRaw("Part Manufacturer> ", false); err != nil {
		return nil, errors.New("Error Reading Component.")
	} else if countI, _ := readUint(countColor("Part Count> "), CountWarning); err != nil {
		return nil, errors.New("Error Reading Component.")
	} else {
		// Success
		return backends.NewComponent(id, countI, name, man), nil
	}
}

// *** Auth Commands
func replLogin(args []string, b backends.Backend) {
	if len(args) == 0 {
		fmt.Println("Please enter your username to login!")
		return
	} else if len(args) > 1 {
		fmt.Println("You didn't enter your password in plaintext did you? You naughty boy.")
		return
	} else {
		user := args[0]
		password := GetPWs()
		cred := backends.NewDummyCredential(user, password, backends.Unknown)
		err := b.GetCredentialManager().Login(cred)
		if err != nil {
			color.Red("Error: " + err.Error())
		} else {
			color.Green("Login Successfull!")
		}
	}
}

func printPermissionTable() {
	fmt.Println("Permission Table:")
	for i := backends.FIRST_CRED; i <= backends.LAST_CRED; i++ {
		fmt.Println(strconv.Itoa(int(i)) + ": " + i.String())
	}
}

func getPermissionIndex() backends.CredentialLevel {
	tmpNum, err := readUint(scanColor("permission index: "), "Please enter a valid permission level")
	var num backends.CredentialLevel
	if err != nil {
		num = backends.DEFAULT_CRED
	} else if tmpNum > uint(int(backends.LAST_CRED)) {
		num = backends.LAST_CRED
	} else {
		num = backends.CredentialLevel(tmpNum)
	}
	return num
}

func replAddUser(args []string, b backends.Backend) {
	if len(args) == 0 {
		fmt.Println("Please enter a user to add.")
		return
	} else if len(args) > 1 {
		fmt.Println("You didn't enter a password in plaintext did you? You naughty boy.")
		return
	} else {
		user := args[0]
		password := GetPWs()

		printPermissionTable()

		num := getPermissionIndex()

		cred := backends.NewDummyCredential(user, password, num)
		if err := b.GetCredentialManager().AddCredential(cred); err != nil {
			color.Red("Error: " + err.Error())
		} else {
			color.Green("User created!")
		}
	}
}

func replDeleteUser(args []string, b backends.Backend) {
	if len(args) == 0 {
		fmt.Println("Please enter a user to delete.")
		return
	} else if len(args) > 1 {
		return
	} else {
		answer, _ := readRaw(authColor("u sure bro? (\"yes\")> "))
		if answer != "yes" {
			color.Red("Aborting!")
			return
		}

		user := args[0]
		cred := backends.NewDummyCredential(user, "", backends.Unknown)
		err := b.GetCredentialManager().RemoveCredential(cred)
		if err != nil {
			color.Red("Error: " + err.Error())
		} else {
			color.Green("User killed!")
		}
	}
}

func replListUsers(args []string, b backends.Backend) {
	if users, err := b.GetCredentialManager().DumpUsers(); err != nil {
		color.Red("Error: " + err.Error())
	} else {
		fmt.Println("Users:")
		for _, name := range users {
			fmt.Println(name.GetUsername() + ": " + name.GetCredentialLevel().String())
		}
	}
}

func replChangeAuth(args []string, b backends.Backend) {
	var user string
	if len(args) == 0 {
		if cred, err := b.GetCredentialManager().CurrentUser(); err != nil {
			fmt.Println("No User logged in!")
			return
		} else {
			user = cred.GetUsername()
		}
	} else {
		user = args[0]
	}
	password := GetPWs()
	cred := backends.NewDummyCredential(user, password, backends.Unknown)
	if err := b.GetCredentialManager().UpdateAuth(cred, password); err != nil {
		color.Red("Error: " + err.Error())
	} else {
		fmt.Println("Changed password of " + user)
	}
}

func replChangePermissionLevel(args []string, b backends.Backend) {
	var user string
	if len(args) == 0 {
		if cred, err := b.GetCredentialManager().CurrentUser(); err != nil {
			fmt.Println("No User logged in!")
			return
		} else {
			user = cred.GetUsername()
		}
	} else {
		user = args[0]
	}
	printPermissionTable()
	perm := getPermissionIndex()
	cred := backends.NewDummyCredential(user, "", perm)
	if err := b.GetCredentialManager().UpdatePermission(cred, perm); err != nil {
		color.Red("Error: " + err.Error())
	} else {
		fmt.Println("Changed permission level of " + user)
	}
}

func replWhoAmI(args []string, b backends.Backend) {
	if len(args) == 0 {
		if user, err := b.GetCredentialManager().CurrentUser(); err != nil {
			color.Red("Error: " + err.Error())
		} else {
			fmt.Println(user.GetUsername() + ": " + user.GetCredentialLevel().String())
		}
	}
}

// *** Misc Commands
func replQuit([]string, backends.Backend) {
	os.Exit(0)
}
