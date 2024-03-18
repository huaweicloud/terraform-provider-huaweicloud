package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type VulScanTaskInfo struct {

	// 任务ID
	Id *string `json:"id,omitempty"`

	// 扫描任务的类型，包含如下：   -manual : 手动扫描任务   -schedule : 定时扫描任务
	ScanType *string `json:"scan_type,omitempty"`

	// 扫描任务开始的时间
	StartTime *int64 `json:"start_time,omitempty"`

	// 扫描任务结束的时间
	EndTime *int64 `json:"end_time,omitempty"`

	// 该任务扫描的漏洞类型列表
	ScanVulTypes *[]string `json:"scan_vul_types,omitempty"`

	// 扫描任务的执行状态，包含如下：   -running : 扫描中   -finished : 扫描完成
	Status *string `json:"status,omitempty"`

	// 该任务处于扫描中状态的主机数量
	ScanningHostNum *int32 `json:"scanning_host_num,omitempty"`

	// 该任务已扫描成功的主机数量
	SuccessHostNum *int32 `json:"success_host_num,omitempty"`

	// 该任务已扫描失败的主机数量
	FailedHostNum *int32 `json:"failed_host_num,omitempty"`
}

func (o VulScanTaskInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VulScanTaskInfo struct{}"
	}

	return strings.Join([]string{"VulScanTaskInfo", string(data)}, " ")
}
