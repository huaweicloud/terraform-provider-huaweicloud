package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type DeleteIkThesaurusRequest struct {

	// 指定要删除自定义词库的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o DeleteIkThesaurusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteIkThesaurusRequest struct{}"
	}

	return strings.Join([]string{"DeleteIkThesaurusRequest", string(data)}, " ")
}
