package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowDomainLocationStatsResponse struct {

	// 数据分组方式
	GroupBy *string `json:"group_by,omitempty"`

	// 按指定的分组方式组织的数据
	Result         map[string]interface{} `json:"result,omitempty"`
	HttpStatusCode int                    `json:"-"`
}

func (o ShowDomainLocationStatsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainLocationStatsResponse struct{}"
	}

	return strings.Join([]string{"ShowDomainLocationStatsResponse", string(data)}, " ")
}
