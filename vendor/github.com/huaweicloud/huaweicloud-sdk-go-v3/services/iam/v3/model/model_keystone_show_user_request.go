package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneShowUserRequest struct {

	// 待查询的IAM用户ID，获取方式请参见：[获取用户ID](https://support.huaweicloud.com/api-iam/iam_17_0002.html)。
	UserId string `json:"user_id"`
}

func (o KeystoneShowUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowUserRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowUserRequest", string(data)}, " ")
}
