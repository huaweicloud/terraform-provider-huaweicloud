package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBucketRegionResponse Response Object
type ShowBucketRegionResponse struct {

	// region ID
	Id *string `json:"id,omitempty"`

	// region 名称
	Name *string `json:"name,omitempty"`

	// 此region是否支持迁移
	Support        *bool `json:"support,omitempty"`
	HttpStatusCode int   `json:"-"`
}

func (o ShowBucketRegionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBucketRegionResponse struct{}"
	}

	return strings.Join([]string{"ShowBucketRegionResponse", string(data)}, " ")
}
