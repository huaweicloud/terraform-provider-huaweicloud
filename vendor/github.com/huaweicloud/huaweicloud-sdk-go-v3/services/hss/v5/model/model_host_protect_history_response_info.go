package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type HostProtectHistoryResponseInfo struct {

	// 静态网页防篡改的检测时间(ms)
	OccrTime *int64 `json:"occr_time,omitempty"`

	// 被篡改文件路径
	FilePath *string `json:"file_path,omitempty"`

	// 文件操作类型   - add: 新增   - delete: 删除   - modify: 修改内容   - attribute: 修改属性   - unknown: 未知
	FileOperation *string `json:"file_operation,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器ip
	HostIp *string `json:"host_ip,omitempty"`

	// 进程ID
	ProcessId *string `json:"process_id,omitempty"`

	// 进程名称
	ProcessName *string `json:"process_name,omitempty"`

	// 进程命令行
	ProcessCmd *string `json:"process_cmd,omitempty"`
}

func (o HostProtectHistoryResponseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "HostProtectHistoryResponseInfo struct{}"
	}

	return strings.Join([]string{"HostProtectHistoryResponseInfo", string(data)}, " ")
}
