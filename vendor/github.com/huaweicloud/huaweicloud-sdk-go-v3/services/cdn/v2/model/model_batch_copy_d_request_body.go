package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCopyDRequestBody 域名复制请求体。
type BatchCopyDRequestBody struct {
	Configs *BatchCopyConfigs `json:"configs"`
}

func (o BatchCopyDRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCopyDRequestBody struct{}"
	}

	return strings.Join([]string{"BatchCopyDRequestBody", string(data)}, " ")
}
