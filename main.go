package main

import (
	"fmt"
	"github.com/nanjingblue/maydb/repl"
	"os"
	"os/user"
	"strings"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	username := currentUser.Username[strings.Index(currentUser.Username, `\`)+1:]
	fmt.Printf("Hello %s!\n", username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
