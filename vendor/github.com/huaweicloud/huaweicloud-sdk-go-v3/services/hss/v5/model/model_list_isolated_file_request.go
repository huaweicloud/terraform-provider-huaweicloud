package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListIsolatedFileRequest Request Object
type ListIsolatedFileRequest struct {

	// Region ID
	Region string `json:"region"`

	// 企业项目ID，查询所有企业项目时填写：all_granted_eps
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 文件路径
	FilePath *string `json:"file_path,omitempty"`

	// 服务器名称
	HostName *string `json:"host_name,omitempty"`

	// 服务器私有IP
	PrivateIp *string `json:"private_ip,omitempty"`

	// 服务器公网IP
	PublicIp *string `json:"public_ip,omitempty"`

	// 文件hash,当前为sha256
	FileHash *string `json:"file_hash,omitempty"`

	// 资产重要性，包含如下3种   - important ：重要资产   - common ：一般资产   - test ：测试资产
	AssetValue *string `json:"asset_value,omitempty"`

	// 偏移量：指定返回记录的开始位置
	Offset *int32 `json:"offset,omitempty"`

	// 每页显示个数
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListIsolatedFileRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListIsolatedFileRequest struct{}"
	}

	return strings.Join([]string{"ListIsolatedFileRequest", string(data)}, " ")
}
