package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"os"
	"smoker/backends"
	"strconv"
	"strings"
)

const INTRO_TEXT = `Welcome to Smoker - The superior beekeper client`
const INTRO_ASCII = `   _________ ___  ____  / /_____  _____              (  )/
  / ___/ __ ` + "`" + `__ \/ __ \/ //_/ _ \/ ___/               )(/
 (__  ) / / / / / /_/ / ,< /  __/ / ________________ ( /)
/____/_/ /_/ /_/\____/_/|_|\___/_/ ()__)____________)))))

`

// TODO make this dummy backend a generic backend
var commands map[string]func([]string, backends.Backend, backends.Credential)

// Main function for smoker
func main() {

	intro()
	initCommands()
	color.Red("WARNING: Using LOCAL Backend. YOU are responsible for ensuring your data is safe!")

	// Make our dummy backend
	backend := backends.NewDummyBackend(10)
	credManager := backend.GetCredentialManager()
	InitialLogin(credManager)

	quit := false
	for !quit {
		if user, err := credManager.CurrentUser(); err != nil {
			fmt.Println(err)
			runCommand(read("> "), backend, user)
		} else {
			runCommand(read("> "), backend, user)
		}
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

const IDWarning string = "was not a valid ID."
const CountWarning string = "was not a valid count."

// Reads uints (ID's and other stuff) interactively from the user in a loop. error on quit
// Pick one of the warnings above or craft your own for the second string value.
func readUint(prompt string, errorMsg string) (uint, error) {
	for {
		idStr, err := readRaw(prompt)
		if err != nil {
			// Clear line due to ctrl-d
			fmt.Println()
			return 0, err
		} else if idStr == "quit" || idStr == "q" {
			return 0, errors.New("User quit")
		} else if len(idStr) == 0 {
			// Blank line, keep going
			continue
		}

		idInt, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			fmt.Println("'" + idStr + "' " + errorMsg)
			continue
		}
		return uint(idInt), nil
	}
}

func readStringLoopRaw(prompt string, allowEmpty bool) (string, error) {
	for {
		result, err := readRaw(prompt)
		if err != nil {
			// Clear line due to ctrl-d
			fmt.Println()
			return "", err
		} else if result == "quit" || result == "q" {
			return "", errors.New("User quit")
		} else if len(result) == 0 && allowEmpty {
			// Blank line, keep going
			continue
		}

		return result, nil
	}
}

func readStringLoop(prompt string) (string, error) {
	return readStringLoopRaw(prompt, true)
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
