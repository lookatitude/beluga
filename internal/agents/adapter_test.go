package agents

import (
	"beluga/pkg/agents"
	"beluga/pkg/agents/adapter"
	"beluga/pkg/orchestration"
	"sync"
	"testing"
	"time"
)

func TestAgentTask(t *testing.T) {
	// Create a test agent
	testAgent := agents.NewDataFetcherAgent("test_task_agent", "test_source", "json")
	config := map[string]interface{}{"test_key": "test_value"}
	if err := testAgent.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}
	
	// Create an agent task
	agentTask := adapter.NewAgentTask(testAgent, "test_task")
	agentTask.WithDependencies("dep1", "dep2")
	
	// Test input setting
	testInput := map[string]string{"key": "value"}
	agentTask.WithInput(testInput)
	
	// Setup result handler
	var resultHandlerCalled bool
	var resultData interface{}
	agentTask.WithResultHandler(func(output interface{}) error {
		resultHandlerCalled = true
		resultData = output
		return nil
	})
	
	// Test ToTask conversion
	task := agentTask.ToTask()
	if task.ID != "test_task" {
		t.Errorf("Expected task ID to be test_task, got %s", task.ID)
	}
	
	if len(task.DependsOn) != 2 || task.DependsOn[0] != "dep1" || task.DependsOn[1] != "dep2" {
		t.Errorf("Dependencies not properly set: %v", task.DependsOn)
	}
	
	// Test execution
	agentTask.SetOutput("test_output") // Mock setting output since our agent doesn't produce real output
	if err := agentTask.Execute(); err != nil {
		t.Errorf("Agent task execution failed: %v", err)
	}
	
	// Check if result handler was called
	if !resultHandlerCalled {
		t.Errorf("Result handler was not called")
	}
	
	if resultData != "test_output" {
		t.Errorf("Expected result data to be 'test_output', got %v", resultData)
	}
}

func TestAgentWorkflow(t *testing.T) {
	// Create a test workflow
	workflow := adapter.NewAgentWorkflow("test_workflow")
	
	// Create agents and tasks
	agent1 := agents.NewDataFetcherAgent("agent1", "source1", "json")
	agent2 := agents.NewAnalyzerAgent("agent2", "test_analysis")
	agent3 := agents.NewDecisionMakerAgent("agent3")
	
	// Initialize agents
	config := map[string]interface{}{"test_key": "test_value"}
	if err := agent1.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent1: %v", err)
	}
	if err := agent2.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent2: %v", err)
	}
	if err := agent3.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent3: %v", err)
	}
	
	// Create agent tasks
	task1 := adapter.NewAgentTask(agent1, "task1")
	task2 := adapter.NewAgentTask(agent2, "task2").WithDependencies("task1")
	task3 := adapter.NewAgentTask(agent3, "task3").WithDependencies("task2")
	
	// Add tasks to workflow
	if err := workflow.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}
	if err := workflow.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}
	if err := workflow.AddTask(task3); err != nil {
		t.Fatalf("Failed to add task3: %v", err)
	}
	
	// Execute workflow sequentially
	if err := workflow.ExecuteSequential(); err != nil {
		t.Errorf("Workflow execution failed: %v", err)
	}
}

func TestAgentMessagingAdapter(t *testing.T) {
	// Create messaging system and adapter
	ms := orchestration.NewMessagingSystem(10)
	adapter := adapter.NewAgentMessagingAdapter(ms)
	
	// Create agents to register
	agent1 := agents.NewDataFetcherAgent("agent1", "source1", "json")
	agent2 := agents.NewAnalyzerAgent("agent2", "test_analysis")
	
	// Initialize agents
	config := map[string]interface{}{"test_key": "test_value"}
	if err := agent1.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent1: %v", err)
	}
	if err := agent2.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent2: %v", err)
	}
	
	// Register agents
	adapter.RegisterAgent("agent1", agent1)
	adapter.RegisterAgent("agent2", agent2)
	
	// Setup message handling
	var wg sync.WaitGroup
	wg.Add(1)
	
	var receivedMessage orchestration.Message
	adapter.RegisterMessageHandler("agent2", func(msg orchestration.Message) error {
		receivedMessage = msg
		wg.Done()
		return nil
	})
	
	// Start message processing
	adapter.StartMessageProcessing()
	
	// Send a message
	payload := map[string]interface{}{"data": "test_data"}
	err := adapter.SendMessage("agent1", "agent2", "test_message", payload)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
	
	// Wait for message to be processed
	waitTimeout := make(chan bool)
	go func() {
		wg.Wait()
		waitTimeout <- false
	}()
	
	// Add a timeout
	go func() {
		time.Sleep(2 * time.Second)
		waitTimeout <- true
	}()
	
	if timedOut := <-waitTimeout; timedOut {
		t.Errorf("Timed out waiting for message to be processed")
	} else {
		// Check the message
		if receivedMessage.Sender != "agent1" {
			t.Errorf("Expected sender to be agent1, got %s", receivedMessage.Sender)
		}
		if receivedMessage.Receiver != "agent2" {
			t.Errorf("Expected receiver to be agent2, got %s", receivedMessage.Receiver)
		}
		if receivedMessage.Type != "test_message" {
			t.Errorf("Expected message type to be test_message, got %s", receivedMessage.Type)
		}
		if receivedMessage.Payload["data"] != "test_data" {
			t.Errorf("Expected payload data to be test_data, got %v", receivedMessage.Payload["data"])
		}
	}
	
	// Stop message processing
	adapter.StopMessageProcessing()
}

func TestAgentIntegration(t *testing.T) {
	// This test demonstrates how all components work together
	// Create agent factory
	factory := agents.NewAgentFactory()
	
	// Create messaging system
	ms := orchestration.NewMessagingSystem(10)
	
	// Create messaging adapter
	msgAdapter := adapter.NewAgentMessagingAdapter(ms)
	
	// Create workflow
	workflow := adapter.NewAgentWorkflow("integrated_workflow")
	
	// Create agents using factory
	dataFetchConfig := map[string]interface{}{
		"data_source": "test_api",
		"data_format": "json",
	}
	
	analyzerConfig := map[string]interface{}{
		"analysis_type": "sentiment",
	}
	
	agent1, err := factory.CreateAgent("DataFetcherAgent", "fetcher", dataFetchConfig)
	if err != nil {
		t.Fatalf("Failed to create data fetcher agent: %v", err)
	}
	
	agent2, err := factory.CreateAgent("AnalyzerAgent", "analyzer", analyzerConfig)
	if err != nil {
		t.Fatalf("Failed to create analyzer agent: %v", err)
	}
	
	// Register agents with messaging adapter
	msgAdapter.RegisterAgent("fetcher", agent1)
	msgAdapter.RegisterAgent("analyzer", agent2)
	
	// Register message handlers
	msgAdapter.RegisterMessageHandler("fetcher", func(msg orchestration.Message) error {
		// Simple handler for test
		return nil
	})
	
	msgAdapter.RegisterMessageHandler("analyzer", func(msg orchestration.Message) error {
		// Simple handler for test
		return nil
	})
	
	// Create tasks for workflow
	task1 := adapter.NewAgentTask(agent1, "fetch_task")
	task2 := adapter.NewAgentTask(agent2, "analyze_task").WithDependencies("fetch_task")
	
	// Add result handler to fetch task to pass data to analyze task
	task1.WithResultHandler(func(output interface{}) error {
		task2.WithInput(output)
		return nil
	})
	
	// Add tasks to workflow
	if err := workflow.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}
	if err := workflow.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}
	
	// Start message processing
	msgAdapter.StartMessageProcessing()
	
	// Execute workflow
	if err := workflow.Execute(); err != nil {
		t.Errorf("Workflow execution failed: %v", err)
	}
	
	// Stop message processing
	msgAdapter.StopMessageProcessing()
}