package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Response Object
type KeystoneShowPermissionResponse struct {
	Role           *RoleResult `json:"role,omitempty"`
	HttpStatusCode int         `json:"-"`
}

func (o KeystoneShowPermissionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowPermissionResponse struct{}"
	}

	return strings.Join([]string{"KeystoneShowPermissionResponse", string(data)}, " ")
}
