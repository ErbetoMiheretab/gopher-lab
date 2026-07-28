package main

import (
	"flag"
	"fmt"
	"os"
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

		err := todos.Edit(*edit, *title, *priority)

		if err != nil {
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
		todos.Print(*filter)

	case *stats:
		// Call our new Map-based stats function
		todos.PrintStats()

	default:
		fmt.Fprintln(os.Stderr, "Invalid command. Use -h for help.")
		os.Exit(1)
	}
}

// Helper function to avoid repeating the save logic and error handling
func saveAndExit(todos *todo.Todos, successMessage string) {
	err := todos.Store(todoFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error saving file:", err)
		os.Exit(1)
	}
	fmt.Println(successMessage)
	os.Exit(0)
}
