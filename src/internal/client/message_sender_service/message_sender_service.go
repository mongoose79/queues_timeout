package message_sender_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	"internal/client/config"
	"log"
	"net/http"
)

func InsertMessages() error {
	for i := 1; i <= config.NumOfQueues; i++ {
		queueName := fmt.Sprintf("queue%d", i)
		writeToQueue(queueName)
	}
	return nil
}

func writeToQueue(queueName string) error {
	for i := 1; i <= config.NumOfMsgsInQueue; i++ {
		msg := fmt.Sprintf("Message %d for queue %s", i, queueName)
		log.Println("Send message " + msg)
		qm := config.QueueMessage{Message: msg}
		jsonValue, _ := json.Marshal(qm)
		url := fmt.Sprintf("%s/create?queue_name=%s", config.BaseUrl, queueName)
		response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			errMsg := fmt.Sprintf("Failed to process request for queue name %s", queueName)
			return stacktrace.Propagate(err, errMsg)
		}
		if response.StatusCode != http.StatusOK {
			return stacktrace.NewError("Error occurred on the server. StatusCode=%d", response.StatusCode)
		}
	}
	return nil
}
