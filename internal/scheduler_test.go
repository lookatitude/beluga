package internal

import (
	"beluga/pkg/orchestration"
	"testing"
)

func TestScheduler(t *testing.T) {
	scheduler := orchestration.NewScheduler()

	task1 := &orchestration.Task{
		ID: "task1",
		Execute: func() error {
			return nil
		},
	}

	task2 := &orchestration.Task{
		ID: "task2",
		DependsOn: []string{"task1"},
		Execute: func() error {
			return nil
		},
	}

	if err := scheduler.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}

	if err := scheduler.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}

	if err := scheduler.Run(); err != nil {
		t.Fatalf("Scheduler run failed: %v", err)
	}
}

func TestSchedulerExecutionStrategies(t *testing.T) {
	scheduler := orchestration.NewScheduler()

	// Define tasks
	task1 := &orchestration.Task{
		ID: "task1",
		Execute: func() error {
			t.Log("Executing task1")
			return nil
		},
	}
	task2 := &orchestration.Task{
		ID: "task2",
		Execute: func() error {
			t.Log("Executing task2")
			return nil
		},
		DependsOn: []string{"task1"},
	}

	// Add tasks to scheduler
	if err := scheduler.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}
	if err := scheduler.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}

	// Test ExecuteSequential
	t.Log("Testing ExecuteSequential")
	if err := scheduler.ExecuteSequential(); err != nil {
		t.Errorf("ExecuteSequential failed: %v", err)
	}

	// Test ExecuteAutonomous
	t.Log("Testing ExecuteAutonomous")
	if err := scheduler.ExecuteAutonomous(); err != nil {
		t.Errorf("ExecuteAutonomous failed: %v", err)
	}
}