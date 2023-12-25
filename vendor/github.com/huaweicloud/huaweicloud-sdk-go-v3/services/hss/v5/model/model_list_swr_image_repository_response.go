package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSwrImageRepositoryResponse Response Object
type ListSwrImageRepositoryResponse struct {

	// 总数
	TotalNum *int32 `json:"total_num,omitempty"`

	// 查询swr镜像仓库镜像列表
	DataList       *[]PrivateImageRepositoryInfo `json:"data_list,omitempty"`
	HttpStatusCode int                           `json:"-"`
}

func (o ListSwrImageRepositoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSwrImageRepositoryResponse struct{}"
	}

	return strings.Join([]string{"ListSwrImageRepositoryResponse", string(data)}, " ")
}
