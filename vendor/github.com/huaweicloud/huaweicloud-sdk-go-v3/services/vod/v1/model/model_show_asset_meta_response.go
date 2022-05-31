package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowAssetMetaResponse struct {

	// 媒资信息列表。
	AssetInfoArray *[]AssetInfo `json:"asset_info_array,omitempty"`

	// 列表是否被截断。  取值如下： - 1：表示本次查询未返回全部结果。 - 0：表示本次查询已经返回了全部结果。
	IsTruncated *int32 `json:"is_truncated,omitempty"`

	// 查询媒资总数。  > 暂只能统计2万个媒资，若您需要查询具体的媒资总数，请[提交工单](https://console.huaweicloud.com/ticket/?#/ticketindex/business?productTypeId=462902cc39a04ab3a429df872021f970)申请。
	Total          *int32 `json:"total,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowAssetMetaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAssetMetaResponse struct{}"
	}

	return strings.Join([]string{"ShowAssetMetaResponse", string(data)}, " ")
}
