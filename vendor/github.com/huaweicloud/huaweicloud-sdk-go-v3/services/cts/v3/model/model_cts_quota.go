package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CtsQuota struct {

	// quota资源类型。
	Type *string `json:"type,omitempty"`

	// 已使用的资源个数。
	Used *int64 `json:"used,omitempty"`

	// 总资源个数。
	Quota *int64 `json:"quota,omitempty"`
}

func (o CtsQuota) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CtsQuota struct{}"
	}

	return strings.Join([]string{"CtsQuota", string(data)}, " ")
}
