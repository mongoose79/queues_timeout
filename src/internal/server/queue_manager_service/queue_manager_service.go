package queue_manager_service

import (
	"fmt"
	"github.com/palantir/stacktrace"
	"internal/server/blockingQueues"
	"sync"
	"time"
)

type QueueManagerService struct {
	Queues map[string]*blockingQueues.BlockingQueue
}

var maxNumOfMsgsInQueue uint64
var timeoutPeriod = 100

var queueManagerServiceInstance *QueueManagerService
var queueManagerServiceOnce sync.Once

func NewQueueManagerService() *QueueManagerService {
	queueManagerServiceOnce.Do(func() {
		queueManagerServiceInstance = &QueueManagerService{
			Queues: make(map[string]*blockingQueues.BlockingQueue),
		}
		maxNumOfMsgsInQueue = 7
	})
	return queueManagerServiceInstance
}

func (q QueueManagerService) ProcessNewMessage(queueName, message string) {
	q.checkAndCreateQueue(queueName)
	q.enqueueMessage(queueName, message)
}

func (q QueueManagerService) checkAndCreateQueue(queueName string) {
	if _, isExist := q.Queues[queueName]; !isExist {
		queue, _ := blockingQueues.NewArrayBlockingQueue(maxNumOfMsgsInQueue)
		q.Queues[queueName] = queue
	}
}

func (q QueueManagerService) enqueueMessage(queueName, message string) {
	q.Queues[queueName].Put(message)
}

func (q QueueManagerService) DequeueMessage(queueName string, timeout int64) (string, error) {
	startTime := time.Now()
	queue := q.Queues[queueName]
	for queue.Size() == 0 {
		time.Sleep(time.Duration(timeoutPeriod) * time.Millisecond)
		duration := time.Now().Sub(startTime)
		if duration.Milliseconds() > timeout {
			return "", stacktrace.NewError("Timeout occurred while trying to dequeue the item from queue %s", queueName)
		}
	}

	if queue.Size() > 0 {
		message, err := q.Queues[queueName].Get()
		if err != nil {
			return "", stacktrace.Propagate(err, "Failed to dequeue the message from %s", queueName)
		}
		return message.(string), nil
	}

	errMsg := fmt.Sprintf("No item in the queue %s", queueName)
	return "", stacktrace.NewError(errMsg)
}
