package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CheckObsBucketsRequest Request Object
type CheckObsBucketsRequest struct {

	// 账户id，参见《云审计服务API参考》“获取账号ID和项目ID”章节。
	DomainId string `json:"domain_id"`

	Body *CheckObsBucketsRequestBody `json:"body,omitempty"`
}

func (o CheckObsBucketsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CheckObsBucketsRequest struct{}"
	}

	return strings.Join([]string{"CheckObsBucketsRequest", string(data)}, " ")
}
