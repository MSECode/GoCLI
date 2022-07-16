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
    NOTE_PATH = "notes"
    TODO_PATH = "todos"
    LIST_PATH = "lists"
    NOTE_TAG  = "NOTE_"
    TODO_TAG  = "TODO_"
    LIST_TAG  = "LIST_"
)

var (
	writeNote string
	readNote  string
    dataDir   string
	notesDir  string
    todosDir  string
    listsDir  string
)

func createDir(path *string) {
    if _, err := os.Stat(*path); os.IsNotExist(err) {
        os.MkdirAll(*path, 0755)
        fmt.Printf("Directory %s created\n", *path)
    } else {
        fmt.Printf("Directory %s already exists\n", *path)
    }
}

func createEnv() {
    fmt.Println("Start instantiating working environment...")

    dataDir = DATA_PATH
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
		fmt.Printf("Directory %s created\n", dataDir)
	} else {
		fmt.Printf("Directory %s already exists\n", dataDir)
    }
    notesDir = filepath.Join(dataDir, NOTE_PATH)
    todosDir = filepath.Join(dataDir, TODO_PATH)
    listsDir = filepath.Join(dataDir, LIST_PATH)
    createDir(&notesDir)
    createDir(&todosDir)
    createDir(&listsDir)
	fmt.Println("Finish instantiating working environment...")
}

func writeTodoToFile(){
    var title   string
    var body    string
    var dueDate string

}

func writeNoteToFile(notesDir *string, currentTime *string) {

    var title string
	var body string

    fmt.Printf("We are writing note to path: %s\n", *notesDir)
	file, err := os.Create(filepath.Join(*notesDir, filepath.Base(NOTE_TAG+*currentTime+FILE_EXT_TEXT)))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buffer := bufio.NewWriter(file)
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

func init() {
	flag.StringVar(&writeNote, "w", "note", "The type of note to be written. Possible types: note, todo, list")
	flag.StringVar(&readNote, "r", "note", "The type of note to be read. Possible types: note, todo, list")
}

func main() {
	flag.Parse()
    if len(os.Args[1:]) < 1 {
        fmt.Println("You must pass a command. For help use -h command")
        os.Exit(1)
    }

	createEnv()
    var currentTime = time.Now().Format(time.RFC3339)
	
    if writeNote == "note" { //default case
	    writeNoteToFile(&notesDir, &currentTime)
    } else if writeNote == "todo" {
    } else if writeNote == "list" {
    } else if writeNote != "" {
        fmt.Println("Only the following values are allowed for -w command: [note, todo, list]")
    }
}
