package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateAssetMetaReq struct {

	// 媒资ID。
	AssetId string `json:"asset_id"`

	// 媒资标题，长度不超过128个字节，UTF-8编码。
	Title *string `json:"title,omitempty"`

	// 媒资描述，长度不超过1024个字节。
	Description *string `json:"description,omitempty"`

	// 媒资分类id。
	CategoryId *int32 `json:"category_id,omitempty"`

	// 媒资标签。  单个标签不超过16个字节，最多不超过16个标签。  多个用逗号分隔，UTF-8编码。
	Tags *string `json:"tags,omitempty"`
}

func (o UpdateAssetMetaReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateAssetMetaReq struct{}"
	}

	return strings.Join([]string{"UpdateAssetMetaReq", string(data)}, " ")
}
