package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateCategoryReq struct {

	// 媒资分类名称，最大64字节。
	Name string `json:"name"`

	// 父分类ID。  若不填，则默认生成一级分类。  根节点分类ID为0。
	ParentId *int32 `json:"parent_id,omitempty"`
}

func (o CreateCategoryReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateCategoryReq struct{}"
	}

	return strings.Join([]string{"CreateCategoryReq", string(data)}, " ")
}
