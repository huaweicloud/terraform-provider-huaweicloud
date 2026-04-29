package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcEipQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_eip_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota"),
				),
			},
		},
	})
}

func testAccDataSourceVpcEipQuotas_basic() string {
	return `
data "huaweicloud_vpc_eip_quotas" "test" {}
`
}
