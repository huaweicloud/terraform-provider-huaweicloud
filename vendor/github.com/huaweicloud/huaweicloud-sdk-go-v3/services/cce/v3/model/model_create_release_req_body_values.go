package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateReleaseReqBodyValues 模板实例的值
type CreateReleaseReqBodyValues struct {

	// 镜像拉取策略
	ImagePullPolicy *string `json:"imagePullPolicy,omitempty"`

	// 镜像标签
	ImageTag *string `json:"imageTag,omitempty"`
}

func (o CreateReleaseReqBodyValues) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateReleaseReqBodyValues struct{}"
	}

	return strings.Join([]string{"CreateReleaseReqBodyValues", string(data)}, " ")
}
