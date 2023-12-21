package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListChartsRequest Request Object
type ListChartsRequest struct {
}

func (o ListChartsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListChartsRequest struct{}"
	}

	return strings.Join([]string{"ListChartsRequest", string(data)}, " ")
}
