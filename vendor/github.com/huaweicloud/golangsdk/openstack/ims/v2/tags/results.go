package tags

import (
	"github.com/huaweicloud/golangsdk"
)

type RespTags struct {
	//contains list of tags, i.e.key value pair
	Tags []Tag `json:"tags"`
}

type commonResult struct {
	golangsdk.Result
}

type ActionResults struct {
	commonResult
}

type GetResult struct {
	commonResult
}

func (r commonResult) Extract() (*RespTags, error) {
	var response RespTags
	err := r.ExtractInto(&response)
	return &response, err
}
