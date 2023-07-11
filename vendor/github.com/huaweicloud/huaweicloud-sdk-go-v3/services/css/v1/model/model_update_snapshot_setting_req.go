package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateSnapshotSettingReq struct {

	// 备份使用的OBS桶的桶名。
	Bucket string `json:"bucket"`

	// 委托名称，委托给CSS，允许CSS调用您的其他云服务。
	Agency string `json:"agency"`

	// 快照在OBS桶中的存放路径。
	BasePath string `json:"basePath"`
}

func (o UpdateSnapshotSettingReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateSnapshotSettingReq struct{}"
	}

	return strings.Join([]string{"UpdateSnapshotSettingReq", string(data)}, " ")
}
