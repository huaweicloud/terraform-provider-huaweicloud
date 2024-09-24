package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePortTagResponse Response Object
type CreatePortTagResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o CreatePortTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePortTagResponse struct{}"
	}

	return strings.Join([]string{"CreatePortTagResponse", string(data)}, " ")
}
