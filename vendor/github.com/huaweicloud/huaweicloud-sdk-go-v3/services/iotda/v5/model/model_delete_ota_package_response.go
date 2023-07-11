package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteOtaPackageResponse Response Object
type DeleteOtaPackageResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteOtaPackageResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteOtaPackageResponse struct{}"
	}

	return strings.Join([]string{"DeleteOtaPackageResponse", string(data)}, " ")
}
