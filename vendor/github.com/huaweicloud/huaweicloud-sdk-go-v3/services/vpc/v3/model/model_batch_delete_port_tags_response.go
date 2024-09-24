package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchDeletePortTagsResponse Response Object
type BatchDeletePortTagsResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchDeletePortTagsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeletePortTagsResponse struct{}"
	}

	return strings.Join([]string{"BatchDeletePortTagsResponse", string(data)}, " ")
}
