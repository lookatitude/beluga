package internal

import (
	"beluga/pkg/orchestration"
	"testing"
	"time"
)

func TestWorkflowMonitor(t *testing.T) {
	monitor := orchestration.NewWorkflowMonitor()
	monitor.StartLogging()

	// Update and verify task states
	monitor.UpdateState("task1", "started")
	if state := monitor.GetState("task1"); state != "started" {
		t.Errorf("Expected state 'started', got '%s'", state)
	}

	monitor.UpdateState("task1", "completed")
	if state := monitor.GetState("task1"); state != "completed" {
		t.Errorf("Expected state 'completed', got '%s'", state)
	}

	// Allow time for logging to process
	time.Sleep(100 * time.Millisecond)
}