package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddHostsGroupResponse Response Object
type AddHostsGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AddHostsGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddHostsGroupResponse struct{}"
	}

	return strings.Join([]string{"AddHostsGroupResponse", string(data)}, " ")
}
