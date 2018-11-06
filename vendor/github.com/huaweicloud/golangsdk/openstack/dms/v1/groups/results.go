package groups

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// GroupCreate response
type GroupCreate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateResult is a struct that contains all the return parameters of creation
type CreateResult struct {
	golangsdk.Result
}

// Extract from CreateResult
func (r CreateResult) Extract() ([]GroupCreate, error) {
	var s struct {
		GroupsCreate []GroupCreate `json:"groups"`
	}
	err := r.Result.ExtractInto(&s)
	return s.GroupsCreate, err
}

// DeleteResult is a struct which contains the result of deletion
type DeleteResult struct {
	golangsdk.ErrResult
}

// Group response
type Group struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	ConsumedMessages     int    `json:"consumed_messages"`
	AvailableMessages    int    `json:"available_messages"`
	ProducedMessages     int    `json:"produced_messages"`
	ProducedDeadletters  int    `json:"produced_deadletters"`
	AvailableDeadletters int    `json:"available_deadletters"`
}

type Groups struct {
	QueueId   string  `json:"queue_id"`
	QueueName string  `json:"queue_name"`
	Details   []Group `json:"groups"`
}

// GroupPage may be embedded in a Page
// that contains all of the results from an operation at once.
type GroupPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a ListResult contains no groups.
func (r GroupPage) IsEmpty() (bool, error) {
	rs, err := ExtractGroups(r)
	return len(rs) == 0, err
}

// ExtractGroups from List
func ExtractGroups(r pagination.Page) ([]Group, error) {
	var s struct {
		Groups []Group `json:"groups"`
	}
	err := (r.(GroupPage)).ExtractInto(&s)
	return s.Groups, err
}
