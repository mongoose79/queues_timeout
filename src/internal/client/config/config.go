package config

type QueueMessage struct {
	Message string `json:"queue_message"`
}

const LogFile = "QueuesTimeout_client.log"

var NumOfQueues = 5
var NumOfMsgsInQueue = 7
var BaseUrl = "http://localhost:8080/api/v1"
