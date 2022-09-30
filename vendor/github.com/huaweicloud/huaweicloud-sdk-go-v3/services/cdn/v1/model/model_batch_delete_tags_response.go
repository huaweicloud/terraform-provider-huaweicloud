package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type BatchDeleteTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeleteTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeleteTagsResponse", string(data)}, " ")
}
