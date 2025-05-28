package cbh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceEcsQuota_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cbh_instance_ecs_quota.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstanceEcsQuota_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("check_status_v6_valid", "true"),
					resource.TestCheckOutput("check_status_valid", "true"),
				),
			},
		},
	})
}

func testDataSourceInstanceEcsQuota_basic() string {
	return `
data "huaweicloud_availability_zones" "test" {
}

data "huaweicloud_cbh_instance_ecs_quota" "test" {
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  resource_spec_code = "cbh.basic.10"
}

output "check_status_v6_valid" {
  value = contains(["sellout", "normal"], data.huaweicloud_cbh_instance_ecs_quota.test.status_v6)
}

output "check_status_valid" {
  value = contains(["sellout", "normal"], data.huaweicloud_cbh_instance_ecs_quota.test.status)
}
`
}
