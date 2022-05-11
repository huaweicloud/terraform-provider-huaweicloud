package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Request Object
type KeystoneShowRegionRequest struct {

	// 待查询的区域ID。可以使用[查询区域列表](https://support.huaweicloud.com/api-iam/iam_05_0001.html)接口获取，控制台获取方法请参见：[获取区域ID](https://console.huaweicloud.com/iam/?agencyId=d15f57bd355d4514bd9618bd648dd432®ion=cn-east-2&locale=zh-cn#/iam/projects)
	RegionId string `json:"region_id"`
}

func (o KeystoneShowRegionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeystoneShowRegionRequest struct{}"
	}

	return strings.Join([]string{"KeystoneShowRegionRequest", string(data)}, " ")
}
