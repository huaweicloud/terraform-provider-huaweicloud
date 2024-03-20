package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTagsRequest Request Object
type ShowTagsRequest struct {

	// 资源id。  > 域名ID
	ResourceId string `json:"resource_id"`
}

func (o ShowTagsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTagsRequest struct{}"
	}

	return strings.Join([]string{"ShowTagsRequest", string(data)}, " ")
}
