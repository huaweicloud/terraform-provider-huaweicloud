package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type ShowProjectDetailsAndStatusResponse struct {
	Project        *ProjectDetailsAndStatusResult `json:"project,omitempty"`
	HttpStatusCode int                            `json:"-"`
}

func (o ShowProjectDetailsAndStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowProjectDetailsAndStatusResponse struct{}"
	}

	return strings.Join([]string{"ShowProjectDetailsAndStatusResponse", string(data)}, " ")
}
