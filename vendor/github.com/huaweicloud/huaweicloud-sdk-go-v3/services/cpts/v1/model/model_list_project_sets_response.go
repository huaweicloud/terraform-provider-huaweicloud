package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ListProjectSetsResponse struct {
	// 状态码

	Code *string `json:"code,omitempty"`
	// 描述

	Message *string `json:"message,omitempty"`
	// 工程集详细信息

	Projects       *[]ProjectsSet `json:"projects,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o ListProjectSetsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProjectSetsResponse struct{}"
	}

	return strings.Join([]string{"ListProjectSetsResponse", string(data)}, " ")
}
