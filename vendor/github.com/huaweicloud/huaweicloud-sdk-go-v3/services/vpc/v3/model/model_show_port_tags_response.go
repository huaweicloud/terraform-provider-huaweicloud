package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPortTagsResponse Response Object
type ShowPortTagsResponse struct {

	// tag对象列表
	Tags *[]ResourceTag `json:"tags,omitempty"`

	// 请求ID
	RequestId      *string `json:"request_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowPortTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPortTagsResponse struct{}"
	}

	return strings.Join([]string{"ShowPortTagsResponse", string(data)}, " ")
}
