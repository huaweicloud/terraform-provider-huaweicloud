package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListPipelinesResponse struct {

	// pipeline列表。
	Pipelines      *[]Pipelines `json:"pipelines,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListPipelinesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPipelinesResponse struct{}"
	}

	return strings.Join([]string{"ListPipelinesResponse", string(data)}, " ")
}
