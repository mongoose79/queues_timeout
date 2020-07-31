package handlers

import (
	"encoding/json"
	"internal/server/queue_manager_service"
	"internal/server/utils"
	"log"
	"net/http"
)

type QueueMessage struct {
	Message string `json:"queue_message"`
}

func ProcessCreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("New message was received from Codota client")
	queueName := r.URL.Query().Get("queue_name")
	if queueName == "" {
		errMsg := "Failed source argument 'queue_name'"
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusBadRequest)
		return
	}

	var qm QueueMessage
	// Try to decode the request body into the string. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&qm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	qms := queue_manager_service.NewQueueManagerService()
	qms.ProcessNewMessage(queueName, qm.Message)
	utils.WriteJSON("StatusOK", w, http.StatusOK)
}
