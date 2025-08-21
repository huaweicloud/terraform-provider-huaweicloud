package rocketmq

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataRocketMQAvailabilityZones_basic(t *testing.T) {
	all := "data.huaweicloud_dms_rocketmq_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRocketMQAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "availability_zones.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.id"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.name"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.code"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.port"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.sold_out"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.resource_availability"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.default_az"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.remain_time"),
					resource.TestCheckResourceAttrSet(all, "availability_zones.0.ipv6_enable"),
				),
			},
		},
	})
}

const testAccDataRocketMQAvailabilityZones_basic = `
data "huaweicloud_dms_rocketmq_availability_zones" "test" {}
`
