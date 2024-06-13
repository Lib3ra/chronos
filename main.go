package main

import (
	"flag"
	"fmt"
)

/*
Todo
create sql db on startup
check if db already created
move db location/filepath

create datamodel for timeentry
add timeentry
track time
edit timeentry
delete timeentry
*/

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
