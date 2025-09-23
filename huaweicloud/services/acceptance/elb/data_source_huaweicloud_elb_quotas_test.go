package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancer"),
					resource.TestCheckResourceAttrSet(dataSource, "listener"),
					resource.TestCheckResourceAttrSet(dataSource, "l7policy"),
					resource.TestCheckResourceAttrSet(dataSource, "pool"),
					resource.TestCheckResourceAttrSet(dataSource, "member"),
					resource.TestCheckResourceAttrSet(dataSource, "security_policy"),
					resource.TestCheckResourceAttrSet(dataSource, "ipgroup"),
					resource.TestCheckResourceAttrSet(dataSource, "ipgroup_max_length"),
					resource.TestCheckResourceAttrSet(dataSource, "healthmonitor"),
					resource.TestCheckResourceAttrSet(dataSource, "certificate"),
					resource.TestCheckResourceAttrSet(dataSource, "ipgroup_bindings"),
					resource.TestCheckResourceAttrSet(dataSource, "listeners_per_loadbalancer"),
					resource.TestCheckResourceAttrSet(dataSource, "listeners_per_pool"),
					resource.TestCheckResourceAttrSet(dataSource, "l7policies_per_listener"),
					resource.TestCheckResourceAttrSet(dataSource, "ipgroups_per_listener"),
					resource.TestCheckResourceAttrSet(dataSource, "members_per_pool"),
					resource.TestCheckResourceAttrSet(dataSource, "pools_per_l7policy"),
					resource.TestCheckResourceAttrSet(dataSource, "condition_per_policy"),
				),
			},
		},
	})
}

func testDataSourceElbQuotas_basic() string {
	return `
data "huaweicloud_elb_quotas" "test" {}
`
}
