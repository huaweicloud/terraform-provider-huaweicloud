package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CommonQueryTaskRsp struct {

	// 任务总数
	Total *int32 `json:"total,omitempty"`
}

func (o CommonQueryTaskRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CommonQueryTaskRsp struct{}"
	}

	return strings.Join([]string{"CommonQueryTaskRsp", string(data)}, " ")
}
