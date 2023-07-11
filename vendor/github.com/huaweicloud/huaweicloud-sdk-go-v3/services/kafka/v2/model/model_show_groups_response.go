package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowGroupsResponse Response Object
type ShowGroupsResponse struct {
	Group          *ShowGroupsRespGroup `json:"group,omitempty"`
	HttpStatusCode int                  `json:"-"`
}

func (o ShowGroupsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowGroupsResponse struct{}"
	}

	return strings.Join([]string{"ShowGroupsResponse", string(data)}, " ")
}
