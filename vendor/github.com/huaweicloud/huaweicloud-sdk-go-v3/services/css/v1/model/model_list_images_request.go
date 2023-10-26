package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListImagesRequest Request Object
type ListImagesRequest struct {

	// 待升级的集群的ID。
	ClusterId string `json:"cluster_id"`

	// 升级目标版本类型： - same：相同版本。 - cross： 跨版本。
	UpgradeType string `json:"upgrade_type"`
}

func (o ListImagesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListImagesRequest struct{}"
	}

	return strings.Join([]string{"ListImagesRequest", string(data)}, " ")
}
