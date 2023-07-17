package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowIkThesaurusRequest Request Object
type ShowIkThesaurusRequest struct {

	// 指定需查询词库状态的集群ID。
	ClusterId string `json:"cluster_id"`
}

func (o ShowIkThesaurusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowIkThesaurusRequest struct{}"
	}

	return strings.Join([]string{"ShowIkThesaurusRequest", string(data)}, " ")
}
