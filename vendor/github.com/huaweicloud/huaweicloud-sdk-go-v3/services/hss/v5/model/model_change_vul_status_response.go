package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ChangeVulStatusResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeVulStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeVulStatusResponse struct{}"
	}

	return strings.Join([]string{"ChangeVulStatusResponse", string(data)}, " ")
}
