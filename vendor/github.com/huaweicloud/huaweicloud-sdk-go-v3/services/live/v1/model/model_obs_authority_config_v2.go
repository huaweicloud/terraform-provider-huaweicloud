package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ObsAuthorityConfigV2 struct {

	// OBS桶名
	Bucket string `json:"bucket"`

	// 操作 - 1：授权 - 0：取消授权
	Operation int32 `json:"operation"`
}

func (o ObsAuthorityConfigV2) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ObsAuthorityConfigV2 struct{}"
	}

	return strings.Join([]string{"ObsAuthorityConfigV2", string(data)}, " ")
}
