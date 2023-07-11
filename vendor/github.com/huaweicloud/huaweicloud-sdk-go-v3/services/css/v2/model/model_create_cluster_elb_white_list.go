package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateClusterElbWhiteList 弹性IP白名单。
type CreateClusterElbWhiteList struct {

	// 是否开启公网访问控制。
	EnableWhiteList bool `json:"enableWhiteList"`

	// 公网访问控制白名单。需要添加白名单的网段或ip，以逗号隔开，不可重复。
	WhiteList *string `json:"whiteList,omitempty"`
}

func (o CreateClusterElbWhiteList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateClusterElbWhiteList struct{}"
	}

	return strings.Join([]string{"CreateClusterElbWhiteList", string(data)}, " ")
}
