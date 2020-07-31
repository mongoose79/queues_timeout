package handlers

import (
	"fmt"
	"internal/server/queue_manager_service"
	"internal/server/utils"
	"log"
	"net/http"
	"strconv"
)

type Data struct {
	Message string `json:"queue_message"`
}

func ProcessReceiveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("New message was received from Codota client...")
	queueName := r.URL.Query().Get("queue_name")
	if queueName == "" {
		errMsg := "Failed source argument 'queue_name'"
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusBadRequest)
		return
	}

	var timeout int64
	timeoutStr := r.URL.Query().Get("timeout")
	if timeoutStr == "" {
		timeout = 10000
	} else {
		var err error
		timeout, err = strconv.ParseInt(timeoutStr, 10, 64)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to parse %s to int", timeoutStr)
			log.Println(errMsg)
		}
	}

	qms := queue_manager_service.NewQueueManagerService()
	message, err := qms.DequeueMessage(queueName, timeout)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to dequeue the message from the queue name %s", queueName)
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusNoContent)
		return
	}
	res := Data{Message: message}
	utils.WriteJSON(res, w, http.StatusOK)
	log.Println("Dequeue message was completed successfully")
}
