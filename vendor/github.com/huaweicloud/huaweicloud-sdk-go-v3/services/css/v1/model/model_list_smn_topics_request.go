package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSmnTopicsRequest Request Object
type ListSmnTopicsRequest struct {

	// 指定待查询的集群ID。
	DomainId string `json:"domain_id"`
}

func (o ListSmnTopicsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSmnTopicsRequest struct{}"
	}

	return strings.Join([]string{"ListSmnTopicsRequest", string(data)}, " ")
}
