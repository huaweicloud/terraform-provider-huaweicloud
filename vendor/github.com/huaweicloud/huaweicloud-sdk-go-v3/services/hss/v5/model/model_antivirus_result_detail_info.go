package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AntivirusResultDetailInfo 病毒查杀结果详情
type AntivirusResultDetailInfo struct {

	// 病毒查杀结果ID
	ResultId *string `json:"result_id,omitempty"`

	// 病毒名称
	MalwareName *string `json:"malware_name,omitempty"`

	// 文件路径
	FilePath *string `json:"file_path,omitempty"`

	// 文件哈希
	FileHash *string `json:"file_hash,omitempty"`

	// 文件大小
	FileSize *int64 `json:"file_size,omitempty"`

	// 文件属主
	FileOwner *string `json:"file_owner,omitempty"`

	// 文件属性
	FileAttr *string `json:"file_attr,omitempty"`

	// 文件创建时间
	FileCtime *int64 `json:"file_ctime,omitempty"`

	// 文件更新时间
	FileMtime *int64 `json:"file_mtime,omitempty"`

	// 更新时间，毫秒
	UpdateTime *int64 `json:"update_time,omitempty"`

	// Agent ID
	AgentId *string `json:"agent_id,omitempty"`
}

func (o AntivirusResultDetailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AntivirusResultDetailInfo struct{}"
	}

	return strings.Join([]string{"AntivirusResultDetailInfo", string(data)}, " ")
}
