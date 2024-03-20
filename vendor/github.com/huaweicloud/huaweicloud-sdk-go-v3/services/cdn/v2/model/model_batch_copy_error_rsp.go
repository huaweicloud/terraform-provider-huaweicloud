package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchCopyErrorRsp 原域名所有配置
type BatchCopyErrorRsp struct {
	Error *BatchCopyErrorRspError `json:"error,omitempty"`
}

func (o BatchCopyErrorRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchCopyErrorRsp struct{}"
	}

	return strings.Join([]string{"BatchCopyErrorRsp", string(data)}, " ")
}
