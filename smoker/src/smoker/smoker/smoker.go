package main

import (
	"bufio"
	"fmt"
	"os"
	"smoker/backends"
	_ "color"
)

const INTRO_TEXT =
	`Welcome to Smoker - The superior beehive client
   _________ ___  ____  / /_____  _____
  / ___/ __ `+"`" + `__ \/ __ \/ //_/ _ \/ ___/
 (__  ) / / / / / /_/ / ,< /  __/ /
/____/_/ /_/ /_/\____/_/|_|\___/_/

for help, type 'help'`

// Main function for smoker
func main() {
	intro()

	dummy := backends.NewDummyBackend(10)
	fmt.Println(dummy)

	component := backends.NewComponent(1, 1, "30 Ω Resistor", "DigiKey")
	dummy.AddComponent(component)
	if c, b, e := dummy.LookupId(1); e != nil {
		fmt.Println("that part was not found")
	} else {
		fmt.Println(c)
		fmt.Println(b)
		// parts := b.GetParts()
		// fmt.Println(parts)
		// fmt.Println(b)
		// fmt.Println(c)
	}

}

func intro() {
	fmt.Println(INTRO_TEXT)
}

// Reads a line of input
func read(s string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(s)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	return text
}
