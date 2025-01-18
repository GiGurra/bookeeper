package common

import (
	"fmt"
	"os"
)

func ExitWithUserError(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
