package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
func (t *Todos) Print() {
	priorityMap := map[int]string{
		1: "High", 2: "Medium", 3: "Low",
	}

	for i, task := range *t {
		status := " "
		if task.Done {
			status = "X"
		}

		priorityStr, exists := priorityMap[task.Priority]

		if !exists {
			priorityStr = "None"
		}

		fmt.Printf("%d. [%s] %s (Priority: %s)\n", i+1, status, task.Title, priorityStr)
	}
}
