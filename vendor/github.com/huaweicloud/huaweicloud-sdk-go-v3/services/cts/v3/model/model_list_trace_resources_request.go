package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListTraceResourcesRequest Request Object
type ListTraceResourcesRequest struct {

	// 账户id，参见《云审计服务API参考》“获取账号ID和项目ID”章节。
	DomainId string `json:"domain_id"`
}

func (o ListTraceResourcesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTraceResourcesRequest struct{}"
	}

	return strings.Join([]string{"ListTraceResourcesRequest", string(data)}, " ")
}
