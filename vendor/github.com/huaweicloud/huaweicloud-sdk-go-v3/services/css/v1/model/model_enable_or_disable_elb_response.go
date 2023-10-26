package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EnableOrDisableElbResponse Response Object
type EnableOrDisableElbResponse struct {

	// 负载均衡器id。
	ElbId          *string `json:"elb_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o EnableOrDisableElbResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnableOrDisableElbResponse struct{}"
	}

	return strings.Join([]string{"EnableOrDisableElbResponse", string(data)}, " ")
}
