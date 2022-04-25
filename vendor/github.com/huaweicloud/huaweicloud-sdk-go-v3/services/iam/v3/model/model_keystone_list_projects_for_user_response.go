package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneListProjectsForUserResponse struct {
	Links *Links `json:"links,omitempty"`

	// 项目信息列表。
	Projects       *[]ProjectResult `json:"projects,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o KeystoneListProjectsForUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneListProjectsForUserResponse struct{}"
	}

	return strings.Join([]string{"KeystoneListProjectsForUserResponse", string(data)}, " ")
}
