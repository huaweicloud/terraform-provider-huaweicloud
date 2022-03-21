package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type UpdateLoginProtectResponse struct {
	LoginProtect   *UpdateLoginProtectRespon `json:"login_protect,omitempty"`
	HttpStatusCode int                       `json:"-"`
}

func (o UpdateLoginProtectResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateLoginProtectResponse struct{}"
	}

	return strings.Join([]string{"UpdateLoginProtectResponse", string(data)}, " ")
}
