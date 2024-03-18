package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateReleaseReqBody 创建模板实例的请求体
type CreateReleaseReqBody struct {

	// 模板ID
	ChartId string `json:"chart_id"`

	// 模板实例描述
	Description *string `json:"description,omitempty"`

	// 模板实例名称
	Name string `json:"name"`

	// 模板实例所在的命名空间
	Namespace string `json:"namespace"`

	// 模板实例版本号
	Version string `json:"version"`

	Parameters *ReleaseReqBodyParams `json:"parameters,omitempty"`

	Values *CreateReleaseReqBodyValues `json:"values"`
}

func (o CreateReleaseReqBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateReleaseReqBody struct{}"
	}

	return strings.Join([]string{"CreateReleaseReqBody", string(data)}, " ")
}
