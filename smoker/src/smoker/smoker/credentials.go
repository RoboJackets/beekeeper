package main

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"smoker/backends"
	"syscall"
)

// Prompts the user to generate a new account or login
func GenerateCredential() backends.Credential {

	user, err := readRaw("Username: ")
	if err != nil {
		fmt.Println("\nSkipping...")
		return backends.NewCleanDummyCredential()
	}

	fmt.Printf("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nSkipping...")
		return backends.NewCleanDummyCredential()
	}
	pass := string(bytePassword)
	fmt.Println()
	cred := backends.NewDummyCredential(user, pass)
	if cred.Verify() {
		color.Green("Login Successful!")
		return cred
	} else {
		color.Red("Login Failed. Using dummy login.")
		return backends.NewCleanDummyCredential()
	}
}
