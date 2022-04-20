package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReportoutlineInfo struct {
	// 平均响应时间

	AvgResponseTime *float64 `json:"avgResponseTime,omitempty"`
	// 分支id

	BranchId *string `json:"branchId,omitempty"`
	// 分支名称

	BranchName *string `json:"branchName,omitempty"`
	// 用例重试次数

	CaseRetry *float64 `json:"caseRetry,omitempty"`
	// 已完成的用例数

	CompleteNum *float64 `json:"completeNum,omitempty"`
	// 持续时间

	Duration *float64 `json:"duration,omitempty"`
	// 结束时间

	EndTime *string `json:"endTime,omitempty"`
	// 已执行用例数

	ExecutedNum *float64 `json:"executedNum,omitempty"`
	// 迭代id

	IterationUri *string `json:"iterationUri,omitempty"`
	// kpi用例数

	KpiCaseCount *float64 `json:"kpiCaseCount,omitempty"`
	// kpi用例执行次数

	KpiCaseExecuteCount *float64 `json:"kpiCaseExecuteCount,omitempty"`
	// kpi用例通过次数

	KpiCasePassCount *float64 `json:"kpiCasePassCount,omitempty"`
	// 最大并发数

	MaxUsers *float64 `json:"maxUsers,omitempty"`
	// 结果为pass的用例数

	PassNum *float64 `json:"passNum,omitempty"`
	// 阶段id

	Stage *float64 `json:"stage,omitempty"`
	// 阶段名称

	StageName *string `json:"stageName,omitempty"`
	// 开始时间

	StartTime *string `json:"startTime,omitempty"`
	// 成功率

	SuccessRate *float64 `json:"successRate,omitempty"`
	// 任务状态

	TaskStatus *float64 `json:"taskStatus,omitempty"`
	// 总用例数

	TotalNum *float64 `json:"totalNum,omitempty"`
	// 性能tps指标

	Tps *float64 `json:"tps,omitempty"`
	// 分支uri

	VersionUri *string `json:"versionUri,omitempty"`
	// 工程id

	ProjectId *string `json:"projectId,omitempty"`
	// 服务id

	ServiceId *string `json:"serviceId,omitempty"`
}

func (o ReportoutlineInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReportoutlineInfo struct{}"
	}

	return strings.Join([]string{"ReportoutlineInfo", string(data)}, " ")
}
