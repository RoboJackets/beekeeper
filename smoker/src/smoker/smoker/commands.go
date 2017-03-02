package main

import (
	"fmt"
	"os"
	"smoker/backends"
	"strconv"
	"strings"
	"text/tabwriter"
)

func initCommands() {

	commands = make(map[string]func([]string, *backends.DummyBackend))
	commands["help"] = replHelp
	commands["h"] = replHelp

	commands["dump"] = replDump
	commands["d"] = replDump

	commands["quit"] = replQuit
	commands["bye"] = replQuit
	commands["q"] = replQuit
}

// Runs a repl command by keying into the command map
func runCommand(prompt string, backend *backends.DummyBackend) {
	prompt = strings.TrimSpace(prompt)
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
	w.Flush()
}

func replDump(s []string, b *backends.DummyBackend) {
	c := b.GetAllComponents()

	if len(c) == 0 {
		fmt.Println("No data is present.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "ID"+"\t"+"NAME"+"\t"+"MANUFACTURER"+"\t")
	for _, v := range c {
		fmt.Fprintln(w, strconv.Itoa(int(v.GetId()))+"\t"+v.GetName()+"\t"+v.GetManufacturer()+"\t")
	}
	w.Flush()
}

func replQuit([]string, *backends.DummyBackend) {
	os.Exit(0)
}
