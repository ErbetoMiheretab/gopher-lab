package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"todocli/todo"
)

const todoFile = ".todos.json"

func main() {

	//Define cli flags

	add := flag.String("add", "", "Task to be included in the Todo list")
	priority := flag.Int("p", 2, "Priority of the task (1=High, 2=Medium, 3=Low)")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	remove := flag.Int("rm", 0, "Item to be removed")

	edit := flag.Int("edit", 0, "Item to be edited")
	title := flag.String("title", "", "New title for the edited task")
	filter := flag.String("filter", "all", "Filter tasks for listing (all, pending, done)")
	stats := flag.Bool("stats", false, "Show statistics about the tasks")

	sortBy := flag.String("sort", "", "Sort tasks by 'priority' or 'date'")

	// Parse the flags
	flag.Parse()

	//Initialize Todo
	todos := &todo.Todos{}

	//Load existing tasks

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//Handle logic based on the flag provided

	switch {
	case *add != "":
		todos.Add(*add, *priority)
		saveAndExit(todos, "Task added successfully")

	case *edit > 0:
		if *title == "" {
			fmt.Fprintln(os.Stderr, "Error: You must provide a new -title when usign -edit")
			os.Exit(1)
		}

		if err := todos.Edit(*edit, *title, *priority); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		saveAndExit(todos, "Task edited successfully!")

	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Println("Task added successfully")
			os.Exit(1)
		}

		saveAndExit(todos, "Task marked as completed!")

	case *remove > 0:
		err := todos.Delete(*remove)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		saveAndExit(todos, "Task removed successfully!")

	case *list:
		// Sort the list before printing if the flag was provided
		if *sortBy != "" {
			todos.Sort(*sortBy)
		}
		todos.Print(*filter)

	case *stats:
		// Call our new Map-based stats function
		todos.PrintStats()

	default:
		fmt.Fprintln(os.Stderr, "Invalid command. Use -h for help.")
		os.Exit(1)
	}
}

// saveAndExit saves the file, triggers a background backup, and exits
func saveAndExit(todos *todo.Todos, successMessage string) {

	if err := todos.Store(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving file:", err)
		os.Exit(1)
	}
	fmt.Println(successMessage)

	// Concurrency
	// A WaitGroup is like a counter. We add 1 to tell Go we have 1 background task running.
	var wg sync.WaitGroup
	wg.Add(1)

	// The 'go' keyword spins this anonymous function off into a background thread.
	go func() {
		defer wg.Done() // When this function finishes, decrement the WaitGroup counter
		createBackup(todoFile)

	}()

	// Wait for the background task to finish before we let the program exit.
	// (If we didn't wait, the program would exit instantly and kill the backup process halfway through).
	wg.Wait()

	os.Exit(0)
}

// createBackup copies the JSON file to a backup file
func createBackup(originalFile string) {
	//OPen original file

	source, err := os.Open(originalFile)

	if err != nil {
		return
	}

	defer source.Close()

	// DEFER: This tells Go "Close this file right before createBackup() finishes".
	// It guarantees the file is closed, even if an error crashes the function earlier.
	// This replaces Python's `with open(..) as source:`

	// Create backup file
	destination, err := os.Create(originalFile + ".bak")
	if err != nil {
		return
	}

	defer destination.Close()

	// Copy data from source to destination
	io.Copy(destination, source)

}
