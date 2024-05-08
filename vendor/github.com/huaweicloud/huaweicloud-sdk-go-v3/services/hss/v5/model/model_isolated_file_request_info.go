package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// IsolatedFileRequestInfo 恢复的文件详情
type IsolatedFileRequestInfo struct {

	// 主机ID
	HostId *string `json:"host_id,omitempty"`

	// 文件哈希
	FileHash *string `json:"file_hash,omitempty"`

	// 文件路径
	FilePath *string `json:"file_path,omitempty"`

	// 文件属性
	FileAttr *string `json:"file_attr,omitempty"`
}

func (o IsolatedFileRequestInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IsolatedFileRequestInfo struct{}"
	}

	return strings.Join([]string{"IsolatedFileRequestInfo", string(data)}, " ")
}
