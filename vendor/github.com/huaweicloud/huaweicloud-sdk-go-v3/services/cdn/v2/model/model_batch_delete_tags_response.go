package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeleteTagsResponse Response Object
type BatchDeleteTagsResponse struct {
	XRequestId     *string `json:"X-Request-Id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BatchDeleteTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteTagsResponse", string(data)}, " ")
}
