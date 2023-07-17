package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteHostsGroupResponse Response Object
type DeleteHostsGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o DeleteHostsGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteHostsGroupResponse struct{}"
	}

	return strings.Join([]string{"DeleteHostsGroupResponse", string(data)}, " ")
}
