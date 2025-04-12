package internal

import (
	"beluga/pkg/orchestration"
	"testing"
)

func TestChainedTask(t *testing.T) {
	// Define tasks
	task1 := &orchestration.ChainedTask{
		ID: "task1",
		Execute: func(input interface{}) (interface{}, error) {
			return input.(int) + 1, nil
		},
	}
	task2 := &orchestration.ChainedTask{
		ID: "task2",
		Execute: func(input interface{}) (interface{}, error) {
			return input.(int) * 2, nil
		},
	}
	task3 := &orchestration.ChainedTask{
		ID: "task3",
		Execute: func(input interface{}) (interface{}, error) {
			return nil, nil
		},
	}

	// Link tasks
	task1.Next = task2
	task2.Next = task3

	// Run the chain
	if err := task1.Run(1); err != nil {
		t.Fatalf("Chained task failed: %v", err)
	}
}