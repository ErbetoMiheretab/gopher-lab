package main



import ("flag"
 "fmt" 
 "os" 
 "todocli/todo"
)


const todoFile = ".todos.json"

func main(){

	//Define cli flags

	add := flag.String("add", "","Task to be included in the Todo list")
	priority := flag.Int("p", 2, "Priority of the task (1=High, 2=Medium, 3=Low)")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	remove := flag.Int("rm", 0, "Item to be removed")


	// Parse the flags
	flag.Parse()


	//Initialize Todo
	todos := &todo.Todos{}

	//Load existing tasks

	if err := todos.Load(todoFile); err !=nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	//Handle logic based on the flag provided

	switch {
	case *add !="":
		todos.Add(*add, *priority)
		err := todos.Store(todoFile)
		if err != nil{
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Task added successfully")

	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil{
			fmt.Println("Task added successfully")
			os.Exit(1)
		}

		err = todos.Store(todoFile)

		if err != nil{
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Task marked as completed!")


	case *remove > 0:
		err := todos.Delete(*remove)
		if err != nil{
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Task removed successfully!")

	case *list:
		todos.Print()


	default:
		fmt.Fprintln(os.Stderr, "Invalid command. Use -h for help.")
		os.Exit(1)
	
		
	}

}