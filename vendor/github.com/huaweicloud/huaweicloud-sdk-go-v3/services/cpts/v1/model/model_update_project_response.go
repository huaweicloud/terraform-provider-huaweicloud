package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateProjectResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateProjectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateProjectResponse struct{}"
	}

	return strings.Join([]string{"UpdateProjectResponse", string(data)}, " ")
}
