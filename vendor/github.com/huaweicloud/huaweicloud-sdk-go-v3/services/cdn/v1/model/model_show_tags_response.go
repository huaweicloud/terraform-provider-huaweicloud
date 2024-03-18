package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTagsResponse Response Object
type ShowTagsResponse struct {

	// 标签列表。
	Tags *[]Map `json:"tags,omitempty"`

	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTagsResponse struct{}"
	}

	return strings.Join([]string{"ShowTagsResponse", string(data)}, " ")
}
