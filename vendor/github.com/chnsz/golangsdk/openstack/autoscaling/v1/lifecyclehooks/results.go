package lifecyclehooks

import "github.com/chnsz/golangsdk"

// Hook is a struct that represents the result of API calling.
type Hook struct {
	// The lifecycle hook name.
	Name string `json:"lifecycle_hook_name"`
	// The lifecycle hook type: INSTANCE_TERMINATING and INSTANCE_LAUNCHING.
	Type string `json:"lifecycle_hook_type"`
	// The default lifecycle hook callback operation: ABANDON and CONTINUE.
	DefaultResult string `json:"default_result"`
	// The lifecycle hook timeout duration in the unit of second.
	Timeout int `json:"default_timeout"`
	// The unique topic in SMN.
	NotificationTopicURN string `json:"notification_topic_urn"`
	// The topic name in SMN.
	NotificationTopicName string `json:"notification_topic_name"`
	// The notification message.
	NotificationMetadata string `json:"notification_metadata"`
	// The UTC-compliant time when the lifecycle hook is created.
	CreateTime string `json:"create_time"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract will deserialize the result to Hook object.
func (r commonResult) Extract() (*Hook, error) {
	var s Hook
	err := r.Result.ExtractInto(&s)
	return &s, err
}

type CreateResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type ListResult struct {
	commonResult
}

// Extract will deserialize the result to Hook array.
func (r ListResult) Extract() (*[]Hook, error) {
	var s struct {
		// An array of lifecycle hooks.
		Hooks []Hook `json:"lifecycle_hooks"`
	}
	err := r.Result.ExtractInto(&s)
	return &s.Hooks, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}
