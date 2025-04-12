package main

import (
	"beluga/pkg/agents"
	"beluga/pkg/agents/adapter"
	"beluga/pkg/agents/config"
	"beluga/pkg/monitoring"
	"beluga/pkg/orchestration"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// Setup logging
	log.Println("Starting Agent Module POC Demo")

	// 1. Load configuration
	configPath := "configs/agents/agents.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cm := config.NewConfigManager()
	err := cm.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// 2. Initialize agent factory
	factory := agents.NewAgentFactory()

	// 3. Create health check manager
	healthManager := monitoring.NewHealthCheckManager()

	// 4. Setup workflow components
	messagingSystem := orchestration.NewMessagingSystem(100)
	messagingAdapter := adapter.NewAgentMessagingAdapter(messagingSystem)
	
	// 5. Create and configure agents
	log.Println("Creating agents...")
	agentConfigs := cm.GetAllAgentConfigs()
	for _, agentConfig := range agentConfigs {
		log.Printf("Creating agent: %s (%s)", agentConfig.Name, agentConfig.Type)
		
		agent, err := factory.CreateAgent(agentConfig.Type, agentConfig.Name, agentConfig.Settings)
		if err != nil {
			log.Printf("Failed to create agent %s: %v", agentConfig.Name, err)
			continue
		}

		// Register with messaging system
		messagingAdapter.RegisterAgent(agentConfig.Name, agent)
		
		// Setup health checks
		healthCheckFunc := monitoring.CreateAgentHealthCheckFunc(func() map[string]interface{} {
			if baseAgent, ok := agent.(*agents.BaseAgent); ok {
				return baseAgent.CheckHealth()
			}
			return map[string]interface{}{
				"name":  agentConfig.Name,
				"state": "unknown",
			}
		})
		
		healthCheck := monitoring.NewHealthCheck(
			"agent_health",
			agentConfig.Name,
			time.Duration(60)*time.Second,
			healthCheckFunc,
		)
		
		healthCheck.RegisterAlert(func(result *monitoring.HealthCheckResult) {
			log.Printf("Health alert for %s: %s - %s", result.ComponentID, result.Status, result.Message)
		})
		
		healthManager.AddCheck(healthCheck)
	}

	// 6. Setup message handling
	setupMessageHandlers(messagingAdapter, factory.Registry)
	
	// 7. Create workflow
	workflow := adapter.NewAgentWorkflow("data_processing_workflow")
	setupWorkflow(workflow, factory.Registry)
	
	// 8. Start services
	log.Println("Starting services...")
	messagingAdapter.StartMessageProcessing()
	healthManager.StartAllChecks()
	
	// 9. Execute workflow
	log.Println("Executing workflow...")
	err = workflow.Execute()
	if err != nil {
		log.Printf("Workflow execution failed: %v", err)
	} else {
		log.Println("Workflow completed successfully")
	}
	
	// 10. Wait for a bit to allow monitoring to run
	log.Println("Workflow completed, monitoring system health...")
	time.Sleep(3 * time.Second)
	
	status, results := healthManager.CheckSystemHealth()
	log.Printf("System health: %s", status)
	
	// 11. Graceful shutdown
	log.Println("Shutting down...")
	healthManager.StopAllChecks()
	messagingAdapter.StopMessageProcessing()
	
	log.Println("Demo completed successfully")
}

// setupMessageHandlers configures message handlers for agents
func setupMessageHandlers(msgAdapter *adapter.AgentMessagingAdapter, registry *agents.AgentRegistry) {
	for _, agentName := range registry.ListAgents() {
		agent, exists := registry.GetAgent(agentName)
		if !exists {
			continue
		}
		
		// Create closure to capture agent name
		currentAgentName := agentName
		msgAdapter.RegisterMessageHandler(currentAgentName, func(msg orchestration.Message) error {
			log.Printf("Agent %s received message from %s: %s", 
				currentAgentName, msg.Sender, msg.Type)
			return nil
		})
	}
}

// setupWorkflow configures the workflow tasks
func setupWorkflow(workflow *adapter.AgentWorkflow, registry *agents.AgentRegistry) {
	// Find our agents
	webFetcher, exists := registry.GetAgent("web_data_fetcher")
	if !exists {
		log.Fatalf("Required agent 'web_data_fetcher' not found")
	}
	
	analyzer, exists := registry.GetAgent("sentiment_analyzer")
	if !exists {
		log.Fatalf("Required agent 'sentiment_analyzer' not found")
	}
	
	recommender, exists := registry.GetAgent("content_recommender")
	if !exists {
		log.Fatalf("Required agent 'content_recommender' not found")
	}
	
	executor, exists := registry.GetAgent("notification_executor")
	if !exists {
		log.Fatalf("Required agent 'notification_executor' not found")
	}
	
	// Create task: data fetching
	fetchTask := adapter.NewAgentTask(webFetcher, "fetch_data")
	fetchTask.WithInput(map[string]interface{}{
		"url": "https://example.com/sample-content",
		"headers": map[string]string{
			"User-Agent": "Beluga Agent",
		},
	})
	
	// Create task: sentiment analysis
	analyzeTask := adapter.NewAgentTask(analyzer, "analyze_sentiment")
	analyzeTask.WithDependencies("fetch_data")
	
	// Create task: recommendation
	recommendTask := adapter.NewAgentTask(recommender, "generate_recommendations")
	recommendTask.WithDependencies("analyze_sentiment")
	
	// Create task: notification
	notifyTask := adapter.NewAgentTask(executor, "send_notifications")
	notifyTask.WithDependencies("generate_recommendations")
	
	// Link tasks with data passing
	fetchTask.WithResultHandler(func(output interface{}) error {
		log.Println("Data fetching completed, passing data to analyzer")
		
		// In a real implementation, we would actually extract real data
		// Here we're using a mock response for demonstration purposes
		mockData := "This is a sample text that will be analyzed for sentiment. " +
			"The sample is very positive and engaging content that users will enjoy."
			
		// Set the mock data as input for the analyzer
		if analyzerAgent, ok := analyzer.(*agents.AnalyzerAgent); ok {
			analyzerAgent.SetInputData(mockData)
		}
		
		return nil
	})
	
	analyzeTask.WithResultHandler(func(output interface{}) error {
		log.Println("Analysis completed, passing results to recommender")
		
		// Set mock analysis data
		mockAnalysis := map[string]interface{}{
			"sentiment": "positive",
			"score": 0.85,
			"keywords": []string{"engaging", "positive", "enjoy"},
		}
		
		// Pass the results to the decision maker
		if decisionAgent, ok := recommender.(*agents.DecisionMakerAgent); ok {
			decisionAgent.SetAnalysisData(mockAnalysis)
		}
		
		return nil
	})
	
	recommendTask.WithResultHandler(func(output interface{}) error {
		log.Println("Recommendation generated, preparing notification")
		
		// Create mock parameters for the executor
		mockParams := map[string]interface{}{
			"recipients": []string{"user@example.com"},
			"subject": "Content Recommendation",
			"template": "recommendation_template",
			"data": map[string]interface{}{
				"recommendations": []string{"Article 1", "Article 2", "Article 3"},
				"userPreferences": true,
			},
		}
		
		// Pass the parameters to the executor
		if executorAgent, ok := executor.(*agents.ExecutorAgent); ok {
			executorAgent.SetParams(mockParams)
		}
		
		return nil
	})
	
	notifyTask.WithResultHandler(func(output interface{}) error {
		log.Println("Notification sent successfully")
		log.Println("Workflow completed")
		return nil
	})
	
	// Add all tasks to workflow
	workflow.AddTask(fetchTask)
	workflow.AddTask(analyzeTask)
	workflow.AddTask(recommendTask)
	workflow.AddTask(notifyTask)
	
	log.Println("Workflow setup completed with 4 tasks")
}