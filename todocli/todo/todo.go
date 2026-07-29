package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"
)

type Task struct {
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	Priority  int        `json:"priority"`
	CreatedAt time.Time  `json:"created_at"`
	Completed *time.Time `json:"completed_at"` // Pointer to time, allows it to be nil (null)

}

type Todos []Task

func (t *Todos) Add(title string, priority int) {
	task := Task{
		Title:     title,
		Done:      false,
		Priority:  priority,
		CreatedAt: time.Now(),
		Completed: nil,
	}

	*t = append(*t, task)
}

func (t *Todos) Complete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("invalid task index")
	}
	ls[index-1].Done = true
	now := time.Now()
	ls[index-1].Completed = &now

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("invalid task index")
	}

	*t = append(ls[:index-1], ls[index:]...)
	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
func (t *Todos) Print(filter string) {

	for i, task := range *t {

		// Filter logic

		if filter == "done" && !task.Done {
			continue //Skip pending tasks
		}

		if filter == "pending" && task.Done {
			continue //Skip completed tasks
		}

		fmt.Printf("%d. %s\n", i+1, task)
	}
}

func (t *Todos) Edit(index int, newTitle string, newPriority int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid task index")
	}

	// Because 'ls' is a reference to the slice, modifying an element
	// directly updates the struct in memory.

	ls[index-1].Title = newTitle
	ls[index-1].Priority = newPriority

	return nil

}

func (t *Todos) PrintStats() {
	// Create maps to store our tallies

	statusCounts := map[string]int{"Pending": 0, "Completed": 0}

	priorityMap := map[int]string{1: "High", 2: "Medium", 3: "Low"}

	priorityCounts := map[string]int{"High": 0, "Medium": 0, "Low": 0, "None": 0}

	// Iterate through the slice and update the maps

	for _, task := range *t {
		if task.Done {
			statusCounts["Completed"]++
		} else {
			statusCounts["Pending"]++
		}

		priorityStr, exists := priorityMap[task.Priority]

		if !exists {
			priorityStr = "None"
		}
		priorityCounts[priorityStr]++

	}

	fmt.Println("--- TODO STATISTICS ---")
	fmt.Printf("Total Tasks: %d\n", len(*t))
	fmt.Printf("Completed: %d | Pending: %d\n", statusCounts["Completed"], statusCounts["Pending"])
	fmt.Println("By Priority:")
	fmt.Printf("  High: %d | Medium: %d | Low: %d\n", priorityCounts["High"], priorityCounts["Medium"], priorityCounts["Low"])

}

func (t Task) String() string {
	status := " "
	if t.Done {
		status = "X"
	}

	priorityMap := map[int]string{1: "High", 2: "Medium", 3: "Low"}
	priorityStr := priorityMap[t.Priority]

	if priorityStr == "" {
		priorityStr = "None"
	}

	dateStr := t.CreatedAt.Format("2006-01-02 15:04")

	return fmt.Sprintf("[%s] %s (Priority: %s , created: %s)", status, t.Title, priorityStr, dateStr)

}

// Sort reorders the tasks in the slice
func (t *Todos) Sort(by string) {
	ls := *t

	// sort.Slice takes the slice, and a function that returns true if item 'i' should come before item 'j'
	sort.Slice(ls, func(i, j int) bool {
		switch by {
		case "priority":
			// Lower number means higher priority (1 = High)
			return ls[i].Priority < ls[j].Priority
		case "date":
			// Older tasks first
			return ls[i].CreatedAt.Before(ls[j].CreatedAt)
		default:
			return false // Don't change order if criteria is unknown
		}
	})

	*t = ls
}
