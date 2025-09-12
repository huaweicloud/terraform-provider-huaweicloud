package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnP2CGatewayJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_p2c_gateway_jobs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnP2CGatewayJobs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.job_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.finished_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.sub_jobs.0.error_message"),

					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceVpnP2CGatewayJobs_basic = `
data "huaweicloud_vpn_p2c_gateway_jobs" "test" {}

locals {
  resource_id = data.huaweicloud_vpn_p2c_gateway_jobs.test.jobs[0].resource_id
}

// filter by resource_id
data "huaweicloud_vpn_p2c_gateway_jobs" "filter_by_resource_id" {
  resource_id = local.resource_id
}

locals {
  filter_result_by_resource_id = [for v in data.huaweicloud_vpn_p2c_gateway_jobs.filter_by_resource_id.jobs[*].resource_id :
    v == local.resource_id]
}

output "is_resource_id_filter_useful" {
  value = length(local.filter_result_by_resource_id) > 0 && alltrue(local.filter_result_by_resource_id)
}
`
