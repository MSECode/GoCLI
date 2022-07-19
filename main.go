package main

import (
	"bufio"
	"errors"
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
    dataDir   string
	notesDir  string
    todosDir  string
    listsDir  string
)

type WriteCommand struct {
    fs *flag.FlagSet

	writeType string
    writeFile string
}

func (w *WriteCommand) Name() string {
    return w.fs.Name()
}

func (w *WriteCommand) Init(args []string) error {
    return w.fs.Parse(args)
}

func (w *WriteCommand) Run() error {

    if w.writeType == "note" { //default case
        fmt.Println("Default options are used")
        fmt.Printf("Writing a: %s to file: %s.\n", strings.ToUpper(w.writeType), w.writeFile)
	    
        writeNoteToFile(&notesDir, w.writeFile)
    } else if w.writeType == "todo" {
        fmt.Printf("Custom options are used: write -t %s -f %s\n", w.writeType, w.writeFile)
        fmt.Printf("Writing a: %s to file: %s.\n", strings.ToUpper(w.writeType), w.writeFile)
        
        writeTodoToFile(&todosDir, w.writeFile)
    } else if w.writeType == "list" {
        fmt.Printf("Custom options are used: write -t %s -f %s\n", w.writeType, w.writeFile)
        fmt.Printf("Writing a: %s to file: %s.\n", strings.ToUpper(w.writeType), w.writeFile)
        
        writeListToFile(&listsDir, w.writeFile)
    } else {
        return errors.New("Invalid option. Usage: write -t [<note>,<todo>,<list>] -f <filename>")
    }

    return nil
}

func NewWriteCommand() *WriteCommand {
    wc := &WriteCommand {
        fs: flag.NewFlagSet("write", flag.ExitOnError),
    }

    wc.fs.StringVar(&wc.writeType, "t", "note", "Type of file to write. Options: [note, todo, list]")
    wc.fs.StringVar(&wc.writeFile, "f", "example", "Name of file to write")

    return wc
}

type ReadCommand struct {
    fs *flag.FlagSet

    readType  string
    readFile  string
}

func (r *ReadCommand) Name() string {
    return r.fs.Name()
}

func (r *ReadCommand) Init(args []string) error {
    return r.fs.Parse(args)
}

func (r *ReadCommand) Run() error {
    fmt.Printf("You have chosen to read a %s.\n", r.readType)

    return nil
}

func NewReadCommand() *ReadCommand {
    rc := &ReadCommand {
        fs: flag.NewFlagSet("read", flag.ExitOnError),
    }

    rc.fs.StringVar(&rc.readType, "t", "note", "Type of file to read. Options: [note, todo, list]")
    rc.fs.StringVar(&rc.readFile, "f", "default", "Name of file to read")

    return rc
}


type Runner interface {
    Init([]string) error
    Run() error
    Name()  string
}

func root(args []string) error {
    if len(args) < 1 {
        return errors.New("You must pass a command. Available: [write, read]")
    }
    
    createEnv()

    cmds := []Runner {
        NewWriteCommand(),
        NewReadCommand(),
    }

    subcommand := os.Args[1]

    for _, cmd := range cmds {
        if cmd.Name() == subcommand {
            cmd.Init(os.Args[2:])
            return cmd.Run()
        }
    }

    return fmt.Errorf("Unknown subcommand: %s", subcommand)
}


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

func writeListToFile(listsDir *string, fileName string) {
    var title   string
    var body    string
    
    fmt.Printf("We are writing a LIST to path: %s\n", *listsDir)
    filePath := filepath.Join(*listsDir, filepath.Base(LIST_TAG+fileName+FILE_EXT_TEXT))
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        file, err := os.Create(filepath.Join(*listsDir, filepath.Base(LIST_TAG+fileName+FILE_EXT_TEXT)))
        buffer := bufio.NewWriter(file)
        if _, err := buffer.WriteString(TITLE_TAG + ": "); err != nil {
            log.Fatalf(ERR_WRITE_MESSAGE, err)
        }
        fmt.Println("Type the title:")
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
        fmt.Printf("Now type the LIST elements separated by %c:\n", ',')
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
        if err != nil {
            panic(err)
        }
        defer file.Close()
    } else {
        fmt.Printf("A LIST with the name %s already exists. Choose another name: \n", filePath)
        var newFileName string
        fmt.Scanf("%s", &newFileName)
        fmt.Printf("New file name: %s\n", newFileName)
        writeListToFile(listsDir, newFileName)
    }
    // Todo: add choise to add elemnt to already created list
}

func writeTodoToFile(todosDir *string, fileName string) {
    var title   string
    var body    string
    var dueDate string
    var years   int
    var months  int
    var days    int

    currentTime := time.Now().Format(time.RFC3339)

    fmt.Printf("We are writing TODO to path: %s\n", *todosDir)
    file, err := os.Create(filepath.Join(*todosDir, filepath.Base(TODO_TAG+fileName+FILE_EXT_TEXT)))
    if err != nil {
        panic(err)
    }
    defer file.Close()

    buffer := bufio.NewWriter(file)
    if _, err := buffer.WriteString(TITLE_TAG + ": "); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Println("Type the title:")
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
    fmt.Println("Now type the body of the TODO:")
    if scanner.Scan() {
        body = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning the body %v", err)
    }
    if _, err := buffer.WriteString(body + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Println("Now type in order years, months and days you need to complete the TODO")
    fmt.Scanf("%d %d %d", &years, &months, &days)
    t, _ := time.Parse(time.RFC3339, currentTime)
    nt := t.AddDate(years, months, days)
    dueDate = nt.Format(time.RFC3339)
    if _, err := buffer.WriteString(dueDate + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    if err := buffer.Flush(); err != nil {
        log.Fatalf("Error while flushing buffer %v", err)
    }

}

func writeNoteToFile(notesDir *string, fileName string) {

    var title string
	var body string

    fmt.Printf("We are writing NOTE to path: %s\n", *notesDir)
	file, err := os.Create(filepath.Join(*notesDir, filepath.Base(NOTE_TAG+fileName+FILE_EXT_TEXT)))
	if err != nil {
		panic(err)
	}
	defer file.Close()
	
    buffer := bufio.NewWriter(file)
    if _, err := buffer.WriteString(TITLE_TAG + ": "); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Type the title of the NOTE:")
    if scanner.Scan() {
        title = scanner.Text()
    }
    if err := scanner.Err(); err != nil {
        log.Fatalf("Error while scanning title %v", err)
    }
    if _, err := buffer.WriteString(title + "\n"); err != nil {
        log.Fatalf(ERR_WRITE_MESSAGE, err)
    }
    fmt.Println("Now type the body of the NOTE:")
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

func main() {
    if err := root(os.Args[1:]); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
//	createEnv()
//    switch os.Args[1] {
//    case "write":
//        writeCommand.Parse(os.Args[2:])    
//    case "read":
//        readCommand.Parse(os.Args[2:])
//    default:
//        fmt.Printf("%q is not a valid command.\n", os.Args[1])
//    }
}
