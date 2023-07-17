package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreateTagsResponse Response Object
type BatchCreateTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchCreateTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreateTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchCreateTagsResponse", string(data)}, " ")
}
