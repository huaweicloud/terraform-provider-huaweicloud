package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateTaskRequestBody CreateTaskRequestBody
type CreateTaskRequestBody struct {

	// 名称
	Name string `json:"name"`

	// 工程id
	ProjectId int32 `json:"project_id"`

	// 事务信息
	Temps *[]string `json:"temps,omitempty"`

	// 压力阶段模式，0：时长模式；1：次数模式；2：混合模式
	OperateMode *int32 `json:"operate_mode,omitempty"`

	// 基准并发
	BenchConcurrent *int32 `json:"bench_concurrent,omitempty"`
}

func (o CreateTaskRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateTaskRequestBody struct{}"
	}

	return strings.Join([]string{"CreateTaskRequestBody", string(data)}, " ")
}
