package global

import "errors"

// Error constants
var (
	ErrorWorkerNotFound     = errors.New("worker not found")
	ErrorTaskNotFound       = errors.New("task not found")
	ErrorAccountNotFound    = errors.New("account not found")
	ErrorNoAvailableWorker  = errors.New("no available worker")
	ErrorTaskAlreadyExist   = errors.New("task already exists")
	ErrorInvalidTaskStatus  = errors.New("invalid task status")
	ErrorInvalidTaskType    = errors.New("invalid task type")
	ErrorQueuePublishFailed = errors.New("failed to publish message to queue")
	ErrorJsonMarshalFailed  = errors.New("failed to marshal JSON")
	ErrorJsonUnmarshalFailed = errors.New("failed to unmarshal JSON")
)
