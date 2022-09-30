package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateSnapshotSettingReq struct {

	// 备份使用的OBS桶的桶名。
	Bucket string `json:"bucket"`

	// 访问OBS使用的IAM委托名称。
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
