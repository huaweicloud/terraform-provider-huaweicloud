package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateSinkTaskQuotaReq struct {

	// 转储任务的总个数。
	SinkMaxTasks int32 `json:"sink_max_tasks"`
}

func (o UpdateSinkTaskQuotaReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSinkTaskQuotaReq struct{}"
	}

	return strings.Join([]string{"UpdateSinkTaskQuotaReq", string(data)}, " ")
}
