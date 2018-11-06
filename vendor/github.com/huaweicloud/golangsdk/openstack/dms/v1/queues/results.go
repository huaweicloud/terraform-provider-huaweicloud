package queues

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// QueueCreate response
type QueueCreate struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	KafkaTopic string `json:"kafka_topic"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() (*QueueCreate, error) {
	var s QueueCreate
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// Queue response
type Queue struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Created          float64 `json:"created"`
	Description      string  `json:"description"`
	QueueMode        string  `json:"queue_mode"`
	Reservation      int     `json:"reservation"`
	MaxMsgSizeByte   int     `json:"max_msg_size_byte"`
	ProducedMessages int     `json:"produced_messages"`
	RedrivePolicy    string  `json:"redrive_policy"`
	MaxConsumeCount  int     `json:"max_consume_count"`
	GroupCount       int     `json:"group_count"`
}

// GetResult contains the body of getting detailed
type GetResult struct {
	golangsdk.Result
}

// Extract from GetResult
func (r GetResult) Extract() (*Queue, error) {
	var s Queue
	err := r.Result.ExtractInto(&s)
	return &s, err
}

// QueuePage may be embedded in a Page
// that contains all of the results from an operation at once.
type QueuePage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no queues.
func (r QueuePage) IsEmpty() (bool, error) {
	rs, err := ExtractQueues(r)
	return len(rs) == 0, err
}

// ExtractQueues from List
func ExtractQueues(r pagination.Page) ([]Queue, error) {
	var s struct {
		Queues []Queue `json:"queues"`
	}
	err := (r.(QueuePage)).ExtractInto(&s)
	return s.Queues, err
}
