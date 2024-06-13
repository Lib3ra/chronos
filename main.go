package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "world", "The name to greet.")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Printf("Hello, %s!\n", *name)
	} else if flag.Arg(0) == "list" {
		fmt.Println("List command")
	} else {
		fmt.Printf("Hello, %s!\n", *name)
	}
}
