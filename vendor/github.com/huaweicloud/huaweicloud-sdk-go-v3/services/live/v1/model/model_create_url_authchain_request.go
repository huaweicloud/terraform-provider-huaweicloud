package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateUrlAuthchainRequest Request Object
type CreateUrlAuthchainRequest struct {
	Body *CreateUrlAuthchainReq `json:"body,omitempty"`
}

func (o CreateUrlAuthchainRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateUrlAuthchainRequest struct{}"
	}

	return strings.Join([]string{"CreateUrlAuthchainRequest", string(data)}, " ")
}
