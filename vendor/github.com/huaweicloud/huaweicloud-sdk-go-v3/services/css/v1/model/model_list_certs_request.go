package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListCertsRequest Request Object
type ListCertsRequest struct {

	// 指定待查询的集群ID。
	ClusterId string `json:"cluster_id"`

	// 指定查询起始值，默认值为1，即从第1个证书开始查询。
	Start *string `json:"start,omitempty"`

	// 指定查询个数，默认值为10，即一次查询10个证书信息。
	Limit *string `json:"limit,omitempty"`

	// 证书类型。defaultCerts为默认证书类型，不指定查询证书类型默认查找自定义证书列表。
	CertsType *string `json:"certsType,omitempty"`
}

func (o ListCertsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListCertsRequest struct{}"
	}

	return strings.Join([]string{"ListCertsRequest", string(data)}, " ")
}
