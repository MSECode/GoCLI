package main

import (
	"flag"
	"fmt"
)

var (
	writeNote string
	readNote  string
)

func init() {
	flag.StringVar(&writeNote, "w", "note", "The type of note to be written. Possible types: note, todo, list")
	flag.StringVar(&readNote, "r", "note", "The type of note to be read. Possible types: note, todo, list")
}

func main() {

	flag.Parse()

	if writeNote != "" { //default case
		var title string
		fmt.Printf("You have choose to write a %s. Type the title:\n", writeNote)
		fmt.Scanln(&title)
		fmt.Printf("This is your title: %s\n", title)
	}
}
