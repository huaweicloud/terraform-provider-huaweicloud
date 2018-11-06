package tags

import (
	"github.com/huaweicloud/golangsdk"
)

type RespTags struct {
	//contains list of tags, i.e.key value pair
	Tags []Tag `json:"tags"`
}

type Resources struct {
	//List of resources i.e. policies
	Resource []Resource `json:"resources"`
	//Total number of resources
	TotalCount int `json:"total_count"`
}

type Resource struct {
	//Resource name
	ResourceName string `json:"resource_name"`
	//Resource ID
	ResourceID string `json:"resource_id"`
	//List of tags
	Tags []Tag `json:"tags"`
}

type commonResult struct {
	golangsdk.Result
}

type CreateResult struct {
	commonResult
}

type DeleteResult struct {
	commonResult
}

type GetResult struct {
	commonResult
}

type QueryResults struct {
	commonResult
}

type ActionResults struct {
	commonResult
}

func (r commonResult) Extract() (*RespTags, error) {
	var response RespTags
	err := r.ExtractInto(&response)
	return &response, err
}

func (r QueryResults) ExtractResources() (*Resources, error) {
	var response Resources
	err := r.ExtractInto(&response)
	return &response, err
}
