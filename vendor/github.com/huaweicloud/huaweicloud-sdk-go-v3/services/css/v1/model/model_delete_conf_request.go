package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteConfRequest struct {

	// 指定删除配置文件的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *DeleteConfReq `json:"body,omitempty"`
}

func (o DeleteConfRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteConfRequest struct{}"
	}

	return strings.Join([]string{"DeleteConfRequest", string(data)}, " ")
}
