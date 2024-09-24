package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCreatePortTagsResponse Response Object
type BatchCreatePortTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchCreatePortTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCreatePortTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchCreatePortTagsResponse", string(data)}, " ")
}
