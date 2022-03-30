package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 指定ID任务的实例信息。
type GetJobInstanceInfoDetail struct {
	// 实例ID。

	Id string `json:"id"`
	// 实例名称。

	Name string `json:"name"`
}

func (o GetJobInstanceInfoDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetJobInstanceInfoDetail struct{}"
	}

	return strings.Join([]string{"GetJobInstanceInfoDetail", string(data)}, " ")
}
