package internal

import (
	"beluga/pkg/orchestration"
	"testing"

	"github.com/trustmaster/goflow"
)

func TestWorkflowGraph(t *testing.T) {
	graph := orchestration.NewWorkflowGraph()

	in := make(chan string, 1)
	out := make(chan string, 1)

	graph.SetInPort("In", in)
	graph.SetOutPort("Out", out)

	// Run the graph as a network
	done := goflow.Run(graph)

	// Send input to the workflow
	in <- "test input"
	close(in)

	// Wait for the network to finish
	<-done

	// Receive output from the workflow
	result := <-out

	expected := "test input processed processed"
	if result != expected {
		t.Fatalf("Expected %s, got %s", expected, result)
	}
}

func TestWorkflowOrchestrator(t *testing.T) {
	// Updated to pass concrete dependencies to `NewWorkflowOrchestrator`
	messagingSystem := orchestration.NewMessagingSystem(10)
	scheduler := orchestration.NewScheduler()
	orchestrator := orchestration.NewWorkflowOrchestrator(messagingSystem, scheduler)

	// Add tasks to the scheduler
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

	if err := orchestrator.Scheduler.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}
	if err := orchestrator.Scheduler.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}

	// Start the workflow orchestrator
	if err := orchestrator.Start(); err != nil {
		t.Errorf("WorkflowOrchestrator failed: %v", err)
	}
}