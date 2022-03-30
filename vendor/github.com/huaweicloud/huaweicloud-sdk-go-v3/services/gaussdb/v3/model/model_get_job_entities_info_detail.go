package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 根据不同的任务，显示不同的内容。
type GetJobEntitiesInfoDetail struct {
}

func (o GetJobEntitiesInfoDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetJobEntitiesInfoDetail struct{}"
	}

	return strings.Join([]string{"GetJobEntitiesInfoDetail", string(data)}, " ")
}
