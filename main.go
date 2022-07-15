package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
    "time"    
)

const (
    FILE_EXT_TEXT = ".txt"
    DATA_PATH = "tmp/data"
    NOTE_PATH = DATA_PATH + "/notes"
    TODO_PATH = DATA_PATH + "/todos"
    LIST_PATH = DATA_PATH + "/lists"
    NOTE_TAG = "NOTE_"
)

var (
	writeNote string
	readNote  string
	dataDir   string
)

func createEnv() {
	dataDir = NOTE_PATH
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
		fmt.Printf("Directory %s created\n", dataDir)
	} else {
		fmt.Printf("Directory %s already exists\n", dataDir)
	}
}

func init() {
	fmt.Println("Start instantiating working environment...")
	createEnv()
	fmt.Println("Finish instantiating working environment...")
	flag.StringVar(&writeNote, "w", "note", "The type of note to be written. Possible types: note, todo, list")
	flag.StringVar(&readNote, "r", "note", "The type of note to be read. Possible types: note, todo, list")
}

func main() {

	flag.Parse()
	var title string
	var body string
    var currentTime = time.Now().Format(time.RFC3339)

    fmt.Printf("The current time is %s\n", currentTime) 
    
	file, err := os.Create(filepath.Join(dataDir, filepath.Base(NOTE_TAG+currentTime+FILE_EXT_TEXT)))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buffer := bufio.NewWriter(file)

	if writeNote != "" { //default case
		fmt.Printf("You have choose to write a %s. Type the title:\n", writeNote)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			title = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if _, err := buffer.WriteString(title + "\n"); err != nil {
			panic(err)
		}
		fmt.Printf("This is your title: %s\n", title)
		fmt.Printf("Now type the body of the %s:\n", writeNote)
		if scanner.Scan() {
			body = scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if _, err := buffer.WriteString(body + "\n"); err != nil {
            panic(err)
		}
		if err := buffer.Flush(); err != nil {
			panic(err)
		}
	}
}
