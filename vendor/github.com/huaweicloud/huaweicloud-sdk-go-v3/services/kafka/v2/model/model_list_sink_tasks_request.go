package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSinkTasksRequest Request Object
type ListSinkTasksRequest struct {

	// 实例转储ID。 请参考[实例生命周期][查询实例]接口返回的数据。
	ConnectorId string `json:"connector_id"`
}

func (o ListSinkTasksRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSinkTasksRequest struct{}"
	}

	return strings.Join([]string{"ListSinkTasksRequest", string(data)}, " ")
}
