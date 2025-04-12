package agents

import (
	"beluga/pkg/agents"
	"beluga/pkg/interfaces"
	"testing"
	"time"
)

func TestBaseAgent(t *testing.T) {
	// Create a new base agent
	agent := agents.NewBaseAgent("test_agent")
	
	// Test initialization
	config := map[string]interface{}{
		"max_retries": 5,
		"retry_delay": 10,
		"test_setting": "value",
	}
	
	if err := agent.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}
	
	// Verify agent state after initialization
	if agent.GetState() != agents.StateReady {
		t.Errorf("Expected agent state to be ready, got %s", agent.GetState())
	}
	
	// Test event system
	eventFired := false
	agent.RegisterEventHandler("state_change", func(data interface{}) error {
		eventFired = true
		state, ok := data.(agents.AgentState)
		if !ok {
			t.Errorf("Event payload is not AgentState type")
		} else if state != agents.StateRunning {
			t.Errorf("Expected state to be running, got %s", state)
		}
		return nil
	})
	
	// Execute agent and check that event was fired
	if err := agent.Execute(); err != nil {
		t.Fatalf("Failed to execute agent: %v", err)
	}
	
	if !eventFired {
		t.Errorf("State change event was not fired")
	}
	
	// Test health check
	health := agent.CheckHealth()
	if health["name"] != "test_agent" {
		t.Errorf("Expected agent name to be test_agent, got %s", health["name"])
	}
	
	// Test graceful shutdown
	if err := agent.GracefulShutdown(time.Millisecond * 100); err != nil {
		t.Errorf("Graceful shutdown failed: %v", err)
	}
	
	if agent.GetState() != agents.StateShutdown {
		t.Errorf("Expected agent state to be shutdown, got %s", agent.GetState())
	}
}

func TestSpecializedAgents(t *testing.T) {
	// Test DataFetcherAgent
	t.Run("DataFetcherAgent", func(t *testing.T) {
		agent := agents.NewDataFetcherAgent("data_fetcher", "test_source", "json")
		if err := agent.Initialize(nil); err == nil {
			t.Errorf("Expected error for nil config, got nil")
		}
		
		config := map[string]interface{}{"key": "value"}
		if err := agent.Initialize(config); err != nil {
			t.Errorf("Failed to initialize agent: %v", err)
		}
		
		if err := agent.Execute(); err != nil {
			t.Errorf("Failed to execute agent: %v", err)
		}
	})
	
	// Test AnalyzerAgent
	t.Run("AnalyzerAgent", func(t *testing.T) {
		agent := agents.NewAnalyzerAgent("analyzer", "test_analysis")
		config := map[string]interface{}{"key": "value"}
		if err := agent.Initialize(config); err != nil {
			t.Errorf("Failed to initialize agent: %v", err)
		}
		
		// Execute should fail with no input data
		if err := agent.Execute(); err == nil {
			t.Errorf("Expected error due to no input data")
		}
		
		// Set input data and try again
		agent.SetInputData("test data")
		if err := agent.Execute(); err != nil {
			t.Errorf("Failed to execute agent: %v", err)
		}
		
		// Check if results were generated
		result := agent.GetAnalysisResult()
		if result == nil {
			t.Errorf("Expected non-nil analysis result")
		}
	})
	
	// Test DecisionMakerAgent
	t.Run("DecisionMakerAgent", func(t *testing.T) {
		agent := agents.NewDecisionMakerAgent("decision_maker")
		config := map[string]interface{}{"key": "value"}
		if err := agent.Initialize(config); err != nil {
			t.Errorf("Failed to initialize agent: %v", err)
		}
		
		// Execute should fail with no analysis data
		if err := agent.Execute(); err == nil {
			t.Errorf("Expected error due to no analysis data")
		}
		
		// Set analysis data and try again
		agent.SetAnalysisData("test analysis")
		if err := agent.Execute(); err != nil {
			t.Errorf("Failed to execute agent: %v", err)
		}
		
		// Check if decision was made
		decision := agent.GetDecision()
		if decision == "" {
			t.Errorf("Expected non-empty decision")
		}
	})
	
	// Test ExecutorAgent
	t.Run("ExecutorAgent", func(t *testing.T) {
		agent := agents.NewExecutorAgent("executor", "test_action", "test_target")
		config := map[string]interface{}{"key": "value"}
		if err := agent.Initialize(config); err != nil {
			t.Errorf("Failed to initialize agent: %v", err)
		}
		
		// Set params and execute
		agent.SetParams(map[string]interface{}{"param1": "value1"})
		if err := agent.Execute(); err != nil {
			t.Errorf("Failed to execute agent: %v", err)
		}
		
		// Check results
		results := agent.GetResults()
		if results == nil {
			t.Errorf("Expected non-nil execution results")
		}
	})
	
	// Test MonitorAgent
	t.Run("MonitorAgent", func(t *testing.T) {
		agent := agents.NewMonitorAgent("monitor", time.Millisecond*100)
		config := map[string]interface{}{"key": "value"}
		if err := agent.Initialize(config); err != nil {
			t.Errorf("Failed to initialize agent: %v", err)
		}
		
		// Add targets to monitor
		agent.AddMonitorTarget("agent1")
		agent.AddMonitorTarget("agent2")
		
		// Execute to start monitoring
		if err := agent.Execute(); err != nil {
			t.Errorf("Failed to execute agent: %v", err)
		}
		
		// Sleep to allow metrics to be collected
		time.Sleep(time.Millisecond * 200)
		
		// Check results
		results := agent.GetMonitorResults()
		if len(results) == 0 {
			t.Errorf("Expected monitoring results")
		}
		
		// Shutdown
		if err := agent.Shutdown(); err != nil {
			t.Errorf("Failed to shutdown agent: %v", err)
		}
	})
}

func TestAgentFactory(t *testing.T) {
	factory := agents.NewAgentFactory()
	
	testCases := []struct {
		name      string
		agentType string
		shouldSucceed bool
	}{
		{"DataFetcherAgent", "DataFetcherAgent", true},
		{"AnalyzerAgent", "AnalyzerAgent", true},
		{"DecisionMakerAgent", "DecisionMakerAgent", true},
		{"ExecutorAgent", "ExecutorAgent", true},
		{"MonitorAgent", "MonitorAgent", true},
		{"InvalidAgent", "InvalidAgent", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := map[string]interface{}{
				"data_source": "test",
				"data_format": "json",
				"analysis_type": "test",
				"action": "test",
				"target": "test",
				"interval_seconds": 60,
			}
			
			agent, err := factory.CreateAgent(tc.agentType, tc.name, config)
			
			if tc.shouldSucceed {
				if err != nil {
					t.Errorf("Expected agent creation to succeed, got error: %v", err)
				}
				if agent == nil {
					t.Errorf("Agent should not be nil")
				}
				
				// Make sure the agent is registered
				_, exists := factory.Registry.GetAgent(tc.name)
				if !exists {
					t.Errorf("Agent was not registered in the registry")
				}
				
				// Test agent interface compliance
				var _ interfaces.Agent = agent
				
			} else {
				if err == nil {
					t.Errorf("Expected agent creation to fail, but it succeeded")
				}
			}
		})
	}
}