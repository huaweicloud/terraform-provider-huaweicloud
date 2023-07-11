package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateResourceTagResponse Response Object
type CreateResourceTagResponse struct {

	// 查询标签下的资源
	FailedResources *[]TagCreateResponseItem `json:"failed_resources,omitempty"`
	HttpStatusCode  int                      `json:"-"`
}

func (o CreateResourceTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateResourceTagResponse struct{}"
	}

	return strings.Join([]string{"CreateResourceTagResponse", string(data)}, " ")
}
