package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

//
type UpdateBucketAuthorizedReq struct {

	// OBS桶名称。
	Bucket string `json:"bucket"`

	// 是否进行桶授权。  取值如下： - 0：取消授权。 - 1：授权。
	Operation string `json:"operation"`
}

func (o UpdateBucketAuthorizedReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateBucketAuthorizedReq struct{}"
	}

	return strings.Join([]string{"UpdateBucketAuthorizedReq", string(data)}, " ")
}
