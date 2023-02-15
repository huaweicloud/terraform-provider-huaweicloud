package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HostProtectHistoryResponseInfo struct {

	// 检测时间
	OccrTime *int64 `json:"occr_time,omitempty"`

	// 被篡改文件路径
	FilePath *string `json:"file_path,omitempty"`

	// 进程ID，操作系统是Windows时返回
	ProcessId *string `json:"process_id,omitempty"`

	// 进程名称，操作系统是Windows时返回
	ProcessName *string `json:"process_name,omitempty"`

	// 进程命令行，操作系统是Windows时返回
	ProcessCmd *string `json:"process_cmd,omitempty"`
}

func (o HostProtectHistoryResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostProtectHistoryResponseInfo struct{}"
	}

	return strings.Join([]string{"HostProtectHistoryResponseInfo", string(data)}, " ")
}
