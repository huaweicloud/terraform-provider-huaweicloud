package stacks

type (
	StackStatus string
	EventType   string
)

var (
	StackStatusCreationComplete     StackStatus = "CREATION_COMPLETE"
	StackStatusDeploymentInProgress StackStatus = "DEPLOYMENT_IN_PROGRESS"
	StackStatusDeploymentFailed     StackStatus = "DEPLOYMENT_FAILED"
	StackStatusDeploymentComplete   StackStatus = "DEPLOYMENT_COMPLETE"
	StackStatusRollbackInProgress   StackStatus = "ROLLBACK_IN_PROGRESS"
	StackStatusRollbackFailed       StackStatus = "ROLLBACK_FAILED"
	StackStatusRollbackComplete     StackStatus = "ROLLBACK_COMPLETE"
	StackStatusDeletionInProgress   StackStatus = "DELETION_IN_PROGRESS"
	StackStatusDeletionFailed       StackStatus = "DELETION_FAILED"

	EventTypeLog                EventType = "LOG"
	EventTypeError              EventType = "ERROR"
	EventTypeDrift              EventType = "DRIFT"
	EventTypeSummary            EventType = "SUMMARY"
	EventTypeCreationInProgress EventType = "CREATION_IN_PROGRESS"
	EventTypeCreationFailed     EventType = "CREATION_FAILED"
	EventTypeCreationComplete   EventType = "CREATION_COMPLETE"
	EventTypeDeletionInProgress EventType = "DELETION_IN_PROGRESS"
	EventTypeDeletionFailed     EventType = "DELETION_FAILED"
	EventTypeDeletionComplete   EventType = "DELETION_COMPLETE"
	EventTypeDeletionSkipped    EventType = "DELETION_SKIPPED"
	EventTypeUpdateInProgress   EventType = "UPDATE_IN_PROGRESS"
	EventTypeUpdateFailed       EventType = "UPDATE_FAILED"
	EventTypeUpdateComplete     EventType = "UPDATE_COMPLETE"
)

// CreateResp is the structure that represents the API response of the 'Create' method, which contains stack ID and the
// deployment ID.
type CreateResp struct {
	// The unique ID of the resource stack.
	StackId string `json:"stack_id"`
	// The unique ID of the resource deployment.
	DeploymentId string `json:"deployment_id"`
}

// listResp is the structure that represents the API response of the 'ListAll' method, which contains the list of stack
// details.
type listResp struct {
	// The list of stack details.
	Stacks []Stack `json:"stacks"`
}

// Stack is the structure that represents the details of the resource stack.
type Stack struct {
	// The name of the resource stack.
	Name string `json:"stack_name"`
	// The description of the resource stack.
	Description string `json:"description"`
	// The unique ID of the resource stack.
	ID string `json:"stack_id"`
	// The current status of the stack.
	// The valid values are as follows:
	// + CREATION_COMPLETE
	// + DEPLOYMENT_IN_PROGRESS
	// + DEPLOYMENT_FAILED
	// + DEPLOYMENT_COMPLETE
	// + ROLLBACK_IN_PROGRESS
	// + ROLLBACK_FAILED
	// + ROLLBACK_COMPLETE
	// + DELETION_IN_PROGRESS
	// + DELETION_FAILED
	Status string `json:"status"`
	// The creation time.
	CreatedAt string `json:"create_time"`
	// The latest update time.
	UpdatedAt string `json:"update_time"`
	// The debug message for current deployment operation.
	StatusMessage string `json:"status_message"`
}

// deployResp is the structure that represents the API response of the 'Deploy' method, which contains the deployment
// ID.
type deployResp struct {
	// The unique ID of the resource deployment.
	DeploymentId string `json:"deployment_id"`
}

// listEventsResp is the structure that represents the API response of the 'ListAllEvents' method, which contains the
// list of the execution events.
type listEventsResp struct {
	// The list of the execution events.
	StackEvents []StackEvent `json:"stack_events"`
}

// StackEvent is the structure that represents the details of the current execution event.
type StackEvent struct {
	// The id name of the resource, that is, the name of the value used by the corresponding resource as the unique id.
	// When the resource is not created, resource_id_key is not returned.
	ResourceIdKey string `json:"resource_id_key"`
	// The id value of the resource, that is, the value used by the corresponding resource as the unique id.
	// When the resource is not created, resource_id_value is not returned.
	ResourceIdValue string `json:"resource_id_value"`
	// The resource name.
	ResourceName string `json:"resource_name"`
	// The resource type.
	ResourceType string `json:"resource_type"`
	// The time when the event occurred, the format is: yyyy-mm-ddTHH:MM:SSZ.
	Time string `json:"time"`
	// The event type.
	// The valid values are as follows:
	// + LOG
	// + ERROR
	// + DRIFT
	// + SUMMARY
	// + CREATION_IN_PROGRESS
	// + CREATION_FAILED
	// + CREATION_COMPLETE
	// + DELETION_IN_PROGRESS
	// + DELETION_FAILED
	// + DELETION_COMPLETE
	// + DELETION_SKIPPED
	// + UPDATE_IN_PROGRESS
	// + UPDATE_FAILED
	// + UPDATE_COMPLETE
	EventType string `json:"event_type"`
	// The message of the current event.
	EventMessage string `json:"event_message"`
	// The time spent changing the resource, in seconds.
	EventSeconds string `json:"event_seconds"`
}
