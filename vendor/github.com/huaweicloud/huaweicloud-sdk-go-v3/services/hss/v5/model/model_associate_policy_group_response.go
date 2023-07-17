package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AssociatePolicyGroupResponse Response Object
type AssociatePolicyGroupResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AssociatePolicyGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssociatePolicyGroupResponse struct{}"
	}

	return strings.Join([]string{"AssociatePolicyGroupResponse", string(data)}, " ")
}
