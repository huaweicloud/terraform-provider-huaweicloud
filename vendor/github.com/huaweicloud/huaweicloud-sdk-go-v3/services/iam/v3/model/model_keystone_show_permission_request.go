package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneShowPermissionRequest struct {

	// 权限ID，获取方式请参见：[获取权限名、权限ID](https://support.huaweicloud.com/api-iam/iam_10_0001.html)。
	RoleId string `json:"role_id"`
}

func (o KeystoneShowPermissionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowPermissionRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowPermissionRequest", string(data)}, " ")
}
