package main

import (
	"encoding/base64"
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

func initCommands() {

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

	commands["rm"] = replRm
	commands["r"] = replRm

	commands["mv"] = replMove
	commands["m"] = replMove

	commands["u"] = replUpdate
	commands["update"] = replUpdate

	commands["moo"] = replMoo

	commands["w"] = replWelcome
	commands["welcome"] = replWelcome

	commands["quit"] = replQuit
	commands["bye"] = replQuit
	commands["q"] = replQuit
}

// Runs a repl command by keying into the command map
func runCommand(prompt string, backend backends.Backend) {
	cmds := strings.Fields(prompt)

	if len(cmds) == 0 {
		return
	}

	if v, p := commands[cmds[0]]; !p {
		fmt.Println("'" + cmds[0] + "' not found!")
	} else {
		v(cmds[1:], backend)
	}
}

func replHelp([]string, backends.Backend) {
	fmt.Println(`Welcome to the Smoker Help Page.

List of Commands:`)
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "(h)elp\tPrints this help message.")
	fmt.Fprintln(w, "(w)elcome\tPrints a overview message, helpful to beginners.")
	fmt.Fprintln(w, "(q)uit\tQuit's the current repl mode. If at top level, quit.")
	fmt.Fprintln(w, "(d)ump\tDumps information for all components.")
	fmt.Fprintln(w, "(r)m <ID> [<ID>, ...]\tDeletes one or more components, by ID.")
	fmt.Fprintln(w, "(m)v <ID> <BIN>\tMoves <ID> to <BIN> if possible.")
	fmt.Fprintln(w, "(u)pdate <ID> <COUNT>\tUpdates the count of ID to COUNT.")
	fmt.Fprintln(w, "(b)ins\tPrints a list of all bins available.")
	fmt.Fprintln(w, "(g)rep <search>\tGreps all information in every component.")
	fmt.Fprintln(w, "(s)can\tLaunches the interactive scanner interface to add/identify parts.\n\t  Takes in Component IDs, which can be printed with a scanner\n\t  (q)uit to exit scanning mode.")
	w.Flush()
}

func printDump(c []backends.Component) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID"+"\t"+"BIN"+"\t"+"NAME"+"\t"+"MANUFACTURER"+"\t"+"COUNT"+"\t")
	for _, v := range c {
		fmt.Fprintln(w, strconv.Itoa(int(v.GetId()))+"\t"+v.GetBin().GetName()+"\t"+v.GetName()+"\t"+v.GetManufacturer()+"\t"+strconv.Itoa(int(v.GetCount()))+"\t")
	}
	w.Flush()

}

func replDump(s []string, b backends.Backend) {
	c := b.GetAllComponents()

	if len(c) == 0 {
		fmt.Println("No data is present.")
	}
	printDump(c)
}

func replBins(s []string, b backends.Backend) {
	allBins := b.GetAllBinNames()
	for i, _ := range allBins {
		fmt.Println(allBins[i])
	}
}

func replScan(s []string, b backends.Backend) {
	for {
		idStr, err := readRaw("Scan an item or enter an ID> ")
		if err != nil {
			fmt.Println()
			break
		} else if idStr == "quit" || idStr == "q" {
			break
		} else if len(idStr) == 0 {
			// Blank line, keep going
			continue
		}

		idInt, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			fmt.Println("'" + idStr + "' was an invalid ID")
			continue
		}

		componentId := uint(idInt)
		_, bin, err := b.LookupId(componentId)
		if err != nil {
			// Add a new item
			comp := genComponent(componentId)
			bin, err = b.AddComponent(comp)
			if err != nil {
				fmt.Println("Failed adding part: " + err.Error())
				continue
			}

			// Display added part
			printBin(bin.GetName())

			for {
				newBin, err := readRaw("Move part?> ")
				if err != nil || len(newBin) == 0 || newBin == "quit" || newBin == "q" {
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
			printBin(bin.GetName())
		}
	}
}

func printBin(s string) {
	c := color.New(color.FgWhite).Add(color.Bold)
	ascii := figlet4go.NewAsciiRender()
	renderStr, _ := ascii.Render(s)
	c.Println(renderStr)
	c.Println(s)
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
	if len(args) != 2 {
		fmt.Println("mv needs exactly 2 arguments.")
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

func replRm(args []string, b backends.Backend) {
	if len(args) == 0 {
		fmt.Println("Input ID values into rm to delete them.")
		return
	}
	components := make([]backends.Component, 0)
	for _, v := range args {
		if id, err := strconv.Atoi(v); err != nil {
			fmt.Println("'" + v + "' is not a valid ID.")
			return
		} else if comp, _, err := b.LookupId(uint(id)); err != nil {
			fmt.Println("'" + v + "' is not a valid ID.")
			return
		} else {
			components = append(components, comp)
		}
	}

	for _, comp := range components {
		if err := b.RemoveComponent(comp); err != nil {
			fmt.Println("[INTERNAL] An internal error occurred when deleting an element. Partial deletion probably occured.")
			return
		}
	}
}

func replUpdate(args []string, b backends.Backend) {
	if len(args) != 2 {
		fmt.Println("Update must take 2 arguments.")
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

func replWelcome(args []string, b backends.Backend) {
	fmt.Println("Welcome to Smoker!\n")
	fmt.Println("Smoker is the CLI frontend to the BeeKeeper inventory suite.")
	fmt.Println("It offers quick and easy access to all functionality beekeeper provides, while staying out of your way as much as possible.")
	fmt.Println()
	fmt.Println("A good place to get started is the 'scan' command, which opens a seperate scanning REPL, allowing you to add new items if available, and print locations of items if found.")
	fmt.Println()
	fmt.Println("For documentation, type 'help'")
}

// Queries the user for the info required to make a component
func genComponent(id uint) backends.Component {
	// Need Count, Name, and Manufacturer
	name, _ := readRaw("Enter Part Name> ")
	man, _ := readRaw("Enter Part Manufacturer> ")
	countI := parseUint("Enter Part Count> ")

	return backends.NewComponent(id, countI, name, man)
}

// Parse uint with internal error checking (loop on fail)
func parseUint(s string) uint {
	var err error = errors.New("dummy error")
	var countI uint
	var i uint64
	for err != nil {
		countS, _ := readRaw(s)
		i, err = strconv.ParseUint(countS, 10, 32)
		countI = uint(i)
		if err != nil {
			fmt.Println("'" + countS + "' is not a valid number.")
		}
	}
	return countI
}

func replQuit([]string, backends.Backend) {
	os.Exit(0)
}
