package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type CreateAssetCategoryResponse struct {

	// 媒资分类名称。
	Name *string `json:"name,omitempty"`

	// 父分类ID。 一级分类父ID为0。
	ParentId *int32 `json:"parent_id,omitempty"`

	// 媒资分类ID。
	Id *int32 `json:"id,omitempty"`

	// 媒资分类层级。  取值如下： - 1：一级分类层级。 - 2：二级分类层级。 - 3：三级分类层级。
	Level *int32 `json:"level,omitempty"`

	// 项目ID。
	ProjectId      *string `json:"projectId,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateAssetCategoryResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateAssetCategoryResponse struct{}"
	}

	return strings.Join([]string{"CreateAssetCategoryResponse", string(data)}, " ")
}
