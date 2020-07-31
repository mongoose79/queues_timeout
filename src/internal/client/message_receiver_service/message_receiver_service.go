package message_receiver_service

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	"internal/client/config"
	"log"
	"net/http"
)

func ReceiveMessages() error {
	for i := 1; i <= config.NumOfQueues; i++ {
		queueName := fmt.Sprintf("queue%d", i)
		log.Println("Start receiving messages from " + queueName)
		err := readFromQueue(queueName)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to receive the message from queue name %s", queueName)
		}
	}
	return nil
}

func readFromQueue(queueName string) error {
	for i := 1; i <= config.NumOfMsgsInQueue; i++ {
		url := fmt.Sprintf("%s/receive?queue_name=%s&timeout=500", config.BaseUrl, queueName)
		response, err := http.Get(url)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to process request for client queue name %s", queueName)
			return stacktrace.Propagate(err, errMsg)
		}

		var qm config.QueueMessage
		err = json.NewDecoder(response.Body).Decode(&qm)
		if err != nil {
			return err
		}
		log.Println(qm.Message)
	}
	return nil
}
