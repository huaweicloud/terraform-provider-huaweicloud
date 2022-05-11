package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type VersionMediatypes struct {

	// 媒体类型。
	Type string `json:"type"`

	// 基础类型。
	Base string `json:"base"`
}

func (o VersionMediatypes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "VersionMediatypes struct{}"
	}

	return strings.Join([]string{"VersionMediatypes", string(data)}, " ")
}
