package internal

import (
	"testing"
	"beluga/pkg/agents"
)

func TestBaseAgentLifecycle(t *testing.T) {
	agent := &agents.BaseAgent{Name: "TestAgent"}
	config := map[string]interface{}{"key": "value"}

	// Test Initialize
	if err := agent.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}

	// Test Execute
	if err := agent.Execute(); err != nil {
		t.Fatalf("Failed to execute agent: %v", err)
	}

	// Test Shutdown
	if err := agent.Shutdown(); err != nil {
		t.Fatalf("Failed to shutdown agent: %v", err)
	}
}

func TestBaseAgentShutdown(t *testing.T) {
	agent := &agents.BaseAgent{Name: "TestAgent"}
	config := map[string]interface{}{"key": "value"}

	// Initialize the agent
	if err := agent.Initialize(config); err != nil {
		t.Fatalf("Failed to initialize agent: %v", err)
	}

	// Shutdown the agent
	if err := agent.Shutdown(); err != nil {
		t.Fatalf("Failed to shutdown agent: %v", err)
	}
}

func TestAgentFactory(t *testing.T) {
	factory := &agents.AgentFactory{}
	config := map[string]interface{}{"key": "value"}

	// Updated to use a valid agent type: DataFetcherAgent
	agent, err := factory.CreateAgent("DataFetcherAgent", "FactoryTestAgent", config)
	if err != nil {
		t.Fatalf("Failed to create agent: %v", err)
	}

	if agent == nil {
		t.Fatal("Agent is nil")
	}

	if err := agent.Execute(); err != nil {
		t.Fatalf("Failed to execute agent: %v", err)
	}
}

func TestAgentIntegration(t *testing.T) {
	t.Run("Agent Initialization and Execution", func(t *testing.T) {
		// Create a mock agent using the factory
		factory := agents.NewAgentFactory()
		agent, err := factory.CreateAgent("DataFetcherAgent", "MockAgent", map[string]interface{}{"key": "value"})
		if err != nil {
			t.Fatalf("Failed to create agent: %v", err)
		}

		// Initialize the agent
		if err := agent.Initialize(map[string]interface{}{"key": "value"}); err != nil {
			t.Errorf("Agent initialization failed: %v", err)
		}

		// Execute the agent
		if err := agent.Execute(); err != nil {
			t.Errorf("Agent execution failed: %v", err)
		}

		// Shutdown the agent
		if err := agent.Shutdown(); err != nil {
			t.Errorf("Agent shutdown failed: %v", err)
		}
	})

	t.Run("Agent Workflow Simulation", func(t *testing.T) {
		// Create agents using the factory
		factory := agents.NewAgentFactory()
		dataFetcher, err := factory.CreateAgent("DataFetcherAgent", "DataFetcher", map[string]interface{}{"source": "API"})
		if err != nil {
			t.Fatalf("Failed to create DataFetcherAgent: %v", err)
		}

		analyzer, err := factory.CreateAgent("AnalyzerAgent", "Analyzer", map[string]interface{}{"algorithm": "ML"})
		if err != nil {
			t.Fatalf("Failed to create AnalyzerAgent: %v", err)
		}

		decisionMaker, err := factory.CreateAgent("DecisionMakerAgent", "DecisionMaker", map[string]interface{}{"rules": "default"})
		if err != nil {
			t.Fatalf("Failed to create DecisionMakerAgent: %v", err)
		}

		// Simulate workflow
		if err := dataFetcher.Execute(); err != nil {
			t.Errorf("DataFetcherAgent execution failed: %v", err)
		}

		if err := analyzer.Execute(); err != nil {
			t.Errorf("AnalyzerAgent execution failed: %v", err)
		}

		if err := decisionMaker.Execute(); err != nil {
			t.Errorf("DecisionMakerAgent execution failed: %v", err)
		}

		// Shutdown agents
		if err := dataFetcher.Shutdown(); err != nil {
			t.Errorf("DataFetcherAgent shutdown failed: %v", err)
		}

		if err := analyzer.Shutdown(); err != nil {
			t.Errorf("AnalyzerAgent shutdown failed: %v", err)
		}

		if err := decisionMaker.Shutdown(); err != nil {
			t.Errorf("DecisionMakerAgent shutdown failed: %v", err)
		}
	})
}