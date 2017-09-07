package main

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"smoker/backends"
	"syscall"
	"os"
)

// steals all your information
func GetPWs() string {
	fmt.Printf("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "";
	} else {
		return string(bytePassword)
	}
}

// Prompts the user to generate a new account or login
func InitialLogin(credManager backends.CredentialManager) {
	user, err := readRaw("Username: ")
	if err != nil {
		fmt.Println("\nSkipping...")
		if r := credManager.Login(backends.NewCleanDummyCredential()); r != nil {
			color.Red("Demo login failed, quitting.")
			os.Exit(1)
		}
	}

	fmt.Printf("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nSkipping...")
		if r := credManager.Login(backends.NewCleanDummyCredential()); r != nil {
			color.Red("Demo login failed, quitting.")
			os.Exit(1)
		}
	}
	pass := string(bytePassword)
	fmt.Println()
	cred := backends.NewDummyCredential(user, pass)
	if r := credManager.Login(cred); r == nil {
		color.Green("Login Successful!")
	} else {
		color.Red("Login Failed. Using dummy login.")
		if r := credManager.Login(backends.NewCleanDummyCredential()); r != nil {
			color.Red("Demo login failed, quitting.")
			os.Exit(1)
		}
	}
}
