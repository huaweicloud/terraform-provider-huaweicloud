package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateProjectResponse Response Object
type CreateProjectResponse struct {

	// 响应码
	Code *string `json:"code,omitempty"`

	// 响应消息
	Message *string `json:"message,omitempty"`

	// 项目ID
	ProjectId      *int32 `json:"project_id,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o CreateProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateProjectResponse struct{}"
	}

	return strings.Join([]string{"CreateProjectResponse", string(data)}, " ")
}
