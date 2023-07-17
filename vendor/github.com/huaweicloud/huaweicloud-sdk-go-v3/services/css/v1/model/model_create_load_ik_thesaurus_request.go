package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateLoadIkThesaurusRequest Request Object
type CreateLoadIkThesaurusRequest struct {

	// 指定配置自定义词库的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *LoadCustomThesaurusReq `json:"body,omitempty"`
}

func (o CreateLoadIkThesaurusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateLoadIkThesaurusRequest struct{}"
	}

	return strings.Join([]string{"CreateLoadIkThesaurusRequest", string(data)}, " ")
}
