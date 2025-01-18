package common

import (
	"fmt"
	"os"
)

func ExitWithUserError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func AskForConfirmation(msg string) bool {
	fmt.Printf("%s [y/N]: ", msg)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	if len(response) == 0 {
		return false
	}
	response = response[:1]
	return response == "y" || response == "Y"
}
