package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// GetUpgradeDetailInfo 升级任务详情。
type GetUpgradeDetailInfo struct {

	// 任务ID。
	Id *string `json:"id,omitempty"`

	// 升级开始时间。
	StartTime *string `json:"startTime,omitempty"`

	// 升级结束时间。
	EndTime *string `json:"endTime,omitempty"`

	// 任务状态。 - RUNNING：升级中。 - SUCCESS：升级成功。 - FAILED：升级失败。 - PARTIAL_FAILED：部分升级失败。
	Status *string `json:"status,omitempty"`

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	AgencyName *string `json:"agencyName,omitempty"`

	ImageInfo *GetTargetImageIdDetail `json:"imageInfo,omitempty"`

	// 所有需要升级的节点名称集合。
	TotalNodes *string `json:"totalNodes,omitempty"`

	// 所有升级完成的节点名称集合。
	CompletedNodes *string `json:"completedNodes,omitempty"`

	// 当前正在升级的节点名称。
	CurrentNodeName *string `json:"currentNodeName,omitempty"`

	// 重试次数。
	ExecuteTimes *string `json:"executeTimes,omitempty"`

	// 集群当前升级的行为。当有query参数时显示该返回值。
	MigrateParam *string `json:"migrateParam,omitempty"`

	// 集群升级预期结果。当有query参数时显示该返回值。
	FinalAzInfoMap *string `json:"finalAzInfoMap,omitempty"`

	CurrentNodeDetail *[]CurrentNodeDetail `json:"currentNodeDetail,omitempty"`
}

func (o GetUpgradeDetailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetUpgradeDetailInfo struct{}"
	}

	return strings.Join([]string{"GetUpgradeDetailInfo", string(data)}, " ")
}
