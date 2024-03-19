package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeVersionInfo 版本信息
type UpgradeVersionInfo struct {

	// 正式版本号，如：v1.19.10
	Release *string `json:"release,omitempty"`

	// 补丁版本号，如r0
	Patch *string `json:"patch,omitempty"`

	// 推荐升级的目标补丁版本号，如r0
	SuggestPatch *string `json:"suggestPatch,omitempty"`

	// 升级目标版本集合
	TargetVersions *[]string `json:"targetVersions,omitempty"`
}

func (o UpgradeVersionInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeVersionInfo struct{}"
	}

	return strings.Join([]string{"UpgradeVersionInfo", string(data)}, " ")
}
