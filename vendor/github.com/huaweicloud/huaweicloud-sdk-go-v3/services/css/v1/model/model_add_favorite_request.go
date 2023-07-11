package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AddFavoriteRequest Request Object
type AddFavoriteRequest struct {

	// 指定添加自定义模板的集群ID。
	ClusterId string `json:"cluster_id"`

	Body *AddFavoriteReq `json:"body,omitempty"`
}

func (o AddFavoriteRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddFavoriteRequest struct{}"
	}

	return strings.Join([]string{"AddFavoriteRequest", string(data)}, " ")
}
