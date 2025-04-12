# Agent Module POC Demo

This folder contains a proof-of-concept demonstration of the Agent Module for the Beluga AI Agent platform. The demo showcases how agents with different roles can work together in a workflow.

## Features Demonstrated

This demo showcases the following features of the Agent module:

1. **Configuration Management** - Loading and using agent configurations from JSON files
2. **Agent Lifecycle** - Agent initialization, execution, and graceful shutdown
3. **Health Monitoring** - Real-time health checks and status reporting
4. **Workflow Orchestration** - Sequential task execution with dependencies
5. **Inter-Agent Communication** - Message passing between agents
6. **Event System** - Event-driven architecture with handlers
7. **Error Handling** - Automatic retries and error reporting

## Workflow Explanation

The demo implements a simple content recommendation workflow:

1. **Data Fetching** - `DataFetcherAgent` retrieves content from a web API
2. **Sentiment Analysis** - `AnalyzerAgent` processes the content and analyzes sentiment
3. **Recommendation** - `DecisionMakerAgent` generates content recommendations
4. **Notification** - `ExecutorAgent` sends notifications with recommendations

## Running the Demo

To run the demo:

```bash
cd /path/to/beluga
go run examples/agent_workflow/main.go
```

You can also specify a custom agent configuration file:

```bash
go run examples/agent_workflow/main.go /path/to/custom/agents.json
```

## Agent Types

The demo includes several specialized agent types:

- **DataFetcherAgent** - Retrieves data from various sources
- **AnalyzerAgent** - Processes and analyzes data to extract insights
- **DecisionMakerAgent** - Makes decisions based on analyzed data
- **ExecutorAgent** - Performs actions based on decisions
- **MonitorAgent** - Monitors system health and performance

## Configuration

Agents are configured using a JSON configuration file located at `configs/agents/agents.json`. The configuration includes:

- Agent type and role
- Agent-specific settings
- Retry policies
- Dependencies between agents
- Logging and monitoring settings

## Architecture

The demo utilizes these key components:

1. **Agent Interface** - Common interface for all agents
2. **BaseAgent** - Base implementation with shared functionality
3. **AgentFactory** - Creates and configures agents dynamically
4. **ConfigManager** - Manages agent configurations
5. **AgentTask** - Adapts agents to work with the task system
6. **AgentWorkflow** - Manages agent task execution
7. **MessagingAdapter** - Facilitates communication between agents
8. **HealthCheckManager** - Monitors agent health and status

## Further Development

This POC demonstrates the core functionality of the Agent module. Future enhancements could include:

- Persistent state management
- AI model integration
- Remote agent execution
- Advanced workflow patterns (parallel, conditional branches)
- Web-based monitoring dashboard
- Agent debugging tools