package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteResourceTagResponse Response Object
type DeleteResourceTagResponse struct {

	// 查询标签下的资源
	FailedResources *[]TagDeleteResponseItem `json:"failed_resources,omitempty"`
	HttpStatusCode  int                      `json:"-"`
}

func (o DeleteResourceTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteResourceTagResponse struct{}"
	}

	return strings.Join([]string{"DeleteResourceTagResponse", string(data)}, " ")
}
