package main

import (
	"color"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"smoker/backends"
	"syscall"
)

type PasswordCredential struct {
	username, password string
}

func (p *PasswordCredential) GetUsername() string {
	return p.username
}
func (p *PasswordCredential) GetAuth() string {
	return p.password
}

func (p *PasswordCredential) Verify() bool {
	// TODO implment for real
	return p.username == "user" && p.password == "password"
}

func NewDummyCredential() backends.Credential {
	return &PasswordCredential{
		username: "user",
		password: "password"}
}

// Prompts the user to generate a new account or login
func GenerateCredential() backends.Credential {

	user, err := readRaw("Username: ")
	if err != nil {
		fmt.Println("\nSkipping...")
		return NewDummyCredential()
	}

	fmt.Printf("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nSkipping...")
		return NewDummyCredential()
	}
	pass := string(bytePassword)
	fmt.Println()
	cred := &PasswordCredential{
		username: user,
		password: pass}
	if cred.Verify() {
		color.Green("Login Successful!")
		return cred
	} else {
		color.Red("Login Failed. Using dummy login.")
		return NewDummyCredential()
	}
}
