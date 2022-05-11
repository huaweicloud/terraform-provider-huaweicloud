package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListMetadataAomPromGetResponse struct {

	// 响应状态。
	Status *string `json:"status,omitempty"`

	// 元数据信息。
	Data           *interface{} `json:"data,omitempty"`
	HttpStatusCode int          `json:"-"`
}

func (o ListMetadataAomPromGetResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMetadataAomPromGetResponse struct{}"
	}

	return strings.Join([]string{"ListMetadataAomPromGetResponse", string(data)}, " ")
}
