package region

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
)

var AP_SOUTHEAST_1 = region.NewRegion("ap-southeast-1", "https://sms.ap-southeast-1.myhuaweicloud.com")

var staticFields = map[string]*region.Region{
	"ap-southeast-1": AP_SOUTHEAST_1,
}

func ValueOf(regionId string) *region.Region {
	if regionId == "" {
		panic("unexpected empty parameter: regionId")
	}
	if _, ok := staticFields[regionId]; ok {
		return staticFields[regionId]
	}
	panic(fmt.Sprintf("unexpected regionId: %s", regionId))
}
