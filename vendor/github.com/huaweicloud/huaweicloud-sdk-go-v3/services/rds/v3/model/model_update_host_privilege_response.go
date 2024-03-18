package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateHostPrivilegeResponse Response Object
type UpdateHostPrivilegeResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateHostPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateHostPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"UpdateHostPrivilegeResponse", string(data)}, " ")
}
