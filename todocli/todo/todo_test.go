package todo

import (
	"testing"
)

// TestAdd tests if a task is correctly added to the list
func TestAdd(t *testing.T) {
	// 1. Setup
	todos := &Todos{}

	// 2. Execution
	todos.Add("Learn Go Testing", 1)

	// 3. Assertions (using standard 'if' statements)
	if len(*todos) != 1 {
		t.Errorf("Expected list length to be 1, got %d", len(*todos))
	}

	if (*todos)[0].Title != "Learn Go Testing" {
		t.Errorf("Expected title 'Learn Go Testing', got '%s'", (*todos)[0].Title)
	}

	if (*todos)[0].Priority != 1 {
		t.Errorf("Expected priority 1, got %d", (*todos)[0].Priority)
	}
}

// TestComplete tests marking a task as done, and checks error handling
func TestComplete(t *testing.T) {
	todos := &Todos{}
	todos.Add("Write Unit Tests", 2)

	// Test a valid completion (Index 1)
	err := todos.Complete(1)
	if err != nil {
		t.Errorf("Expected no error when completing task 1, got: %s", err)
	}

	if (*todos)[0].Done != true {
		t.Errorf("Expected task to be marked as Done (true), got %v", (*todos)[0].Done)
	}

	// Test an invalid completion (Index 99 - doesn't exist)
	err = todos.Complete(99)
	if err == nil {
		t.Errorf("Expected an error when completing a non-existent task, but got nil")
	}
}

// TestDelete tests removing a task
func TestDelete(t *testing.T) {
	todos := &Todos{}
	todos.Add("Task 1", 1)
	todos.Add("Task 2", 2)
	todos.Add("Task 3", 3)

	// Delete the middle task (Index 2)
	err := todos.Delete(2)
	if err != nil {
		t.Errorf("Expected no error, got: %s", err)
	}

	// Length should now be 2
	if len(*todos) != 2 {
		t.Fatalf("Expected list length to be 2, got %d", len(*todos)) // Fatalf stops the test immediately
	}

	// The new second task (Index 2 in CLI, index 1 in slice) should now be "Task 3"
	if (*todos)[1].Title != "Task 3" {
		t.Errorf("Expected remaining task to be 'Task 3', got '%s'", (*todos)[1].Title)
	}
}
