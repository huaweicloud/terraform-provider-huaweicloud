package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 配置文件信息。
type UpdateYmlsReqEdit struct {
	Modify *UpdateYmlsReqEditModify `json:"modify"`
}

func (o UpdateYmlsReqEdit) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateYmlsReqEdit struct{}"
	}

	return strings.Join([]string{"UpdateYmlsReqEdit", string(data)}, " ")
}
