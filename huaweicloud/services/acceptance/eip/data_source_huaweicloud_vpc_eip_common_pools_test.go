package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcEipCommonPools_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_eip_common_pools.all"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcEipCommonPools_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.allow_share_bandwidth_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "common_pools.0.available"),
				),
			},
		},
	})
}

const testDataSourceVpcEipCommonPools_basic = `data "huaweicloud_vpc_eip_common_pools" "all" {}`
