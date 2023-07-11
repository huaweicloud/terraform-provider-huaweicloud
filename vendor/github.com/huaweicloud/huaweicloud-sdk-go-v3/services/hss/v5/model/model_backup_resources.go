package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BackupResources 开启备份功能新版参数，必填；若为空代表兼容之前绑定HSS_projectid的存储库
type BackupResources struct {

	// 选择需要绑定的存储库ID，不为空
	VaultId *string `json:"vault_id,omitempty"`

	// 需要开启备份功能的主机情况列表
	ResourceList *[]ResourceInfo `json:"resource_list,omitempty"`
}

func (o BackupResources) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupResources struct{}"
	}

	return strings.Join([]string{"BackupResources", string(data)}, " ")
}
