package main

import (
	"bufio"
	"fmt"
	"os"
	"smoker/dummybackend"
)

// Main function for smoker
func main() {
	fmt.Println("hello world")
	x := dummybackend.Add(1, 2)
	fmt.Println(x)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
