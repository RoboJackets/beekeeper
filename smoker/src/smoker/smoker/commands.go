package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"smoker/backends"
	"text/tabwriter"
)

func initCommands() {

	commands = make(map[string]func([]string, *backends.DummyBackend))
	commands["help"] = replHelp
	commands["h"] = replHelp

	commands["dump"] = replDump
	commands["d"] = replDump

	commands["scan"] = replScan
	commands["s"] = replScan

	commands["quit"] = replQuit
	commands["bye"] = replQuit
	commands["q"] = replQuit
}

// Runs a repl command by keying into the command map
func runCommand(prompt string, backend *backends.DummyBackend) {
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

func replHelp([]string, *backends.DummyBackend) {
	fmt.Println(`Welcome to the Smoker Help Page.

List of Commands:`)
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "(h)elp\tPrints this help message")
	fmt.Fprintln(w, "(d)ump\tDumps information for all components")
	fmt.Fprintln(w, "(s)can\tLaunches the interactive scanner interface to add/identify parts")
	w.Flush()
}

func replDump(s []string, b *backends.DummyBackend) {
	c := b.GetAllComponents()

	if len(c) == 0 {
		fmt.Println("No data is present.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID"+"\t"+"BIN"+"\t"+"NAME"+"\t"+"MANUFACTURER"+"\t")
	for _, v := range c {
		fmt.Fprintln(w, strconv.Itoa(int(v.GetId()))+"\t"+v.GetBin().GetName()+"\t"+v.GetName()+"\t"+v.GetManufacturer()+"\t")
	}
	w.Flush()
}

func replScan (s []string, b *backends.DummyBackend) {
	for {
		idStr, err := readRaw("Scan an item or enter an ID> ")
		if err != nil {
			fmt.Println()
			break
		} else if (idStr == "quit") {
			break
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
			b.AddComponent(comp)
		} else {
			// Display found item
			fmt.Println(bin.GetName())
		}
	}
}

// Queries the user for the info required to make a component
func genComponent(id uint) *backends.Component {
	// Need Count, Name, and Manufacturer
	name, _ := readRaw("Enter Part Name> ")
	countI := parseUint("Enter Part Count> ")
	man, _ := readRaw("Enter Part Manufacturer> ")

	return backends.NewComponent(id, countI, name, man)
}

// Parse uint with internal error checking (loop on fail)
func parseUint(s string) uint {
	var err error
	var countI uint
	for err != nil {
		countS, _ := readRaw(s)
		i, err := strconv.ParseUint(countS, 10, 32)
		countI = uint(i)
		if err != nil {
			fmt.Println("'" + countS + "' is not a valid number.")
		}
	}
	return countI
}

func replQuit([]string, *backends.DummyBackend) {
	os.Exit(0)
}
