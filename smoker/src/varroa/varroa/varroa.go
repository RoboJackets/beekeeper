package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"smoker/backends"
)

func main() {
	backend := backends.NewDummyBackend(10)
	credManager := backend.GetCredentialManager()
	if r := credManager.Login(backends.NewCleanDummyCredential()); r != nil {
		color.Red("Demo login failed, quitting.")
		os.Exit(1)
	}

	if user, err := credManager.CurrentUser(); err != nil {
		fmt.Println(err)
	} else {
		startUi(backend, credManager, user)
	}
}
