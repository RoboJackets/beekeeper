package main

import (
	"bufio"
	"github.com/fatih/color"
	"fmt"
	"os"
	"smoker/backends"
	"strings"
)

const INTRO_TEXT = `Welcome to Smoker - The superior beekeper client`
const INTRO_ASCII = `   _________ ___  ____  / /_____  _____
  / ___/ __ ` + "`" + `__ \/ __ \/ //_/ _ \/ ___/
 (__  ) / / / / / /_/ / ,< /  __/ /
/____/_/ /_/ /_/\____/_/|_|\___/_/

`

// TODO make this dummy backend a generic backend
var commands map[string]func([]string, backends.Backend)

// Main function for smoker
func main() {

	intro()
	initCommands()
	color.Red("WARNING: Using DUMMY Backend. Your data will not be stored.")

	// Get auth credentials
	auth := GenerateCredential()

	// Make our dummy backend
	backend := backends.NewDummyBackend(auth, 10)

	quit := false
	for !quit {
		runCommand(read("> "), backend)
	}
}

func intro() {
	fmt.Println(INTRO_TEXT)
	color.Cyan(INTRO_ASCII)

	c := color.New(color.FgRed).Add(color.Bold)
	fmt.Print("For help, type '")
	c.Printf("help")
	fmt.Println("' in the REPL")

	fmt.Println()
	color.Yellow("This software is in heavy development. Please report bugs to RoboJackets/beekeeper")
}

func readRaw(s string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(s)
	text, err := reader.ReadString('\n')
	return strings.TrimSpace(text), err
}

// Reads a line of input
func read(s string) string {
	text, err := readRaw(s)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return text
}
