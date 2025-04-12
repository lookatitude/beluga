package internal

import (
	"beluga/pkg/orchestration"
	"testing"
	"time"
)

func TestMessagingSystem(t *testing.T) {
	ms := orchestration.NewMessagingSystem(10)

	msg := orchestration.Message{
		ID:        "1",
		Timestamp: time.Now(),
		Sender:    "AgentA",
		Receiver:  "AgentB",
		Type:      "task_request",
		Payload:   map[string]interface{}{"key": "value"},
	}

	// Test SendMessage
	if err := ms.SendMessage(msg); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Test ReceiveMessage
	receivedMsg, err := ms.ReceiveMessage()
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	if receivedMsg.ID != msg.ID {
		t.Fatalf("Expected message ID %s, got %s", msg.ID, receivedMsg.ID)
	}

	// Test ValidateMessage
	if err := orchestration.ValidateMessage(msg); err != nil {
		t.Fatalf("Message validation failed: %v", err)
	}

	// Test SerializeMessage
	serialized, err := orchestration.SerializeMessage(msg)
	if err != nil {
		t.Fatalf("Failed to serialize message: %v", err)
	}

	// Test DeserializeMessage
	deserialized, err := orchestration.DeserializeMessage(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize message: %v", err)
	}

	if deserialized.ID != msg.ID {
		t.Fatalf("Expected deserialized message ID %s, got %s", msg.ID, deserialized.ID)
	}
}

func TestSendMessageWithRetry(t *testing.T) {
	ms := orchestration.NewMessagingSystem(1)
	msg := orchestration.Message{
		ID:        "1",
		Timestamp: time.Now(),
		Sender:    "AgentA",
		Receiver:  "AgentB",
		Type:      "task_request",
		Payload:   map[string]interface{}{"key": "value"},
	}

	err := ms.SendMessageWithRetry(msg, 3, 100*time.Millisecond)
	if err != nil {
		t.Errorf("SendMessageWithRetry failed: %v", err)
	}

	receivedMsg, err := ms.ReceiveMessage()
	if err != nil {
		t.Errorf("Failed to receive message: %v", err)
	}

	if receivedMsg.ID != msg.ID {
		t.Errorf("Expected message ID %s, got %s", msg.ID, receivedMsg.ID)
	}
}