package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSchedules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_schedules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the test data in advance under the CFW firewall.
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSchedules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.schedule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ref_count"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSchedules_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cfw_schedules" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}

# Filter by name
locals {
  name = data.huaweicloud_cfw_schedules.test.records[0].name
}

data "huaweicloud_cfw_schedules" "filter_by_name" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name      = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_cfw_schedules.filter_by_name.records[*].name :
    v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, testAccDatasourceFirewalls_basic())
}
