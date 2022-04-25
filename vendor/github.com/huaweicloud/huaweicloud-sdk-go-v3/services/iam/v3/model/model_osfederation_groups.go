package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type OsfederationGroups struct {

	// 用户组ID。
	Id string `json:"id"`

	// 用户组名称。
	Name string `json:"name"`
}

func (o OsfederationGroups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OsfederationGroups struct{}"
	}

	return strings.Join([]string{"OsfederationGroups", string(data)}, " ")
}
