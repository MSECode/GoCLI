package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
    TITLE_TAG = "TITLE"
    ERR_WRITE_MESSAGE = "Error while writing to file %v"
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

func writeListToFile(listsDir *string, currentTime *string) {
    var title   string
    var body    string
    
    fmt.Printf("We are writing %s to path: %s\n", writeNote, *listsDir)
    file, err := os.Create(filepath.Join(*listsDir, filepath.Base(LIST_TAG+*currentTime+FILE_EXT_TEXT)))
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    buffer := bufio.NewWriter(file)
    if _, err := buffer.WriteString(TITLE_TAG + ": "); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("You have chosen to write a %s. Type the title: \n", writeNote)
    scanner := bufio.NewScanner(os.Stdin)
    if scanner.Scan() {
        title = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning title %v", err)
    }
    if _, err := buffer.WriteString(title + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("Now type the %s elements separated by %c:\n", writeNote, ',')
    if scanner.Scan() {
        body = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning the body %v", err)
    }
    listEl := strings.Split(body, ",")
    for _, el := range listEl {
        if _, err := buffer.WriteString(el + "\n"); err != nil {
            log.Fatalf(ERR_WRITE_MESSAGE, err)
        }
    }
    if err := buffer.Flush(); err != nil {
        log.Fatalf("Error while flushing buffer %v", err)
    }
}

func writeTodoToFile(todosDir *string, currentTime *string) {
    var title   string
    var body    string
    var dueDate string
    var years   int
    var months  int
    var days    int

    fmt.Printf("We are writing todo to path: %s\n", *todosDir)
    file, err := os.Create(filepath.Join(*todosDir, filepath.Base(TODO_TAG+*currentTime+FILE_EXT_TEXT)))
    if err != nil {
        panic(err)
    }
    defer file.Close()

    buffer := bufio.NewWriter(file)
    if _, err := buffer.WriteString(TITLE_TAG + ": "); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("You have chosen to write a %s. Type the title: \n", writeNote)
    scanner := bufio.NewScanner(os.Stdin)
    if scanner.Scan() {
        title = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning title %v", err)
    }
    if _, err := buffer.WriteString(title + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("Now type the body of the %s:\n", writeNote)
    if scanner.Scan() {
        body = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning the body %v", err)
    }
    if _, err := buffer.WriteString(body + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("Now type in order years, months and days you need to complete the %s:\n", writeNote)
    fmt.Scanf("%d %d %d", &years, &months, &days)
    t, _ := time.Parse(time.RFC3339, *currentTime)
    nt := t.AddDate(years, months, days)
    dueDate = nt.Format(time.RFC3339)
    if _, err := buffer.WriteString(dueDate + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    if err := buffer.Flush(); err != nil {
        log.Fatalf("Error while flushing buffer %v", err)
    }

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
    if _, err := buffer.WriteString(TITLE_TAG + ": "); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("You have chosen to write a %s. Type the title:\n", writeNote)
    scanner := bufio.NewScanner(os.Stdin)
    if scanner.Scan() {
        title = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning title %v", err)
    }
    if _, err := buffer.WriteString(title + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Printf("Now type the body of the %s:\n", writeNote)
    if scanner.Scan() {
        body = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning the body %v", err)
    }
    if _, err := buffer.WriteString(body + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    if err := buffer.Flush(); err != nil {
        log.Fatalf("Error while flushing buffer %v", err)
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
        writeTodoToFile(&todosDir, &currentTime)
    } else if writeNote == "list" {
        writeListToFile(&listsDir, &currentTime)
    } else if writeNote != "" {
        fmt.Println("Only the following values are allowed for -w command: [note, todo, list]")
    }
}
