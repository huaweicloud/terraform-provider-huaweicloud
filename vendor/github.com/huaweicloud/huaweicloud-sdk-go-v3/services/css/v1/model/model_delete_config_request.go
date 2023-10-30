package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteConfigRequest Request Object
type DeleteConfigRequest struct {

	// 指定删除配置文件的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *DeleteConfReq `json:"body,omitempty"`
}

func (o DeleteConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConfigRequest struct{}"
	}

	return strings.Join([]string{"DeleteConfigRequest", string(data)}, " ")
}
