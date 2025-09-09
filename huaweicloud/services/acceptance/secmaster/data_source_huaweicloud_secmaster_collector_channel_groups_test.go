package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCollectorChannelGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_collector_channel_groups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a collector channel group.
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCollectorChannelGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.name"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCollectorChannelGroups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_collector_channel_groups" "test" {
  workspace_id = "%[1]s"
}

locals {
  name = data.huaweicloud_secmaster_collector_channel_groups.test.groups[0].name
}

data "huaweicloud_secmaster_collector_channel_groups" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_collector_channel_groups.filter_by_name.groups) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_collector_channel_groups.filter_by_name.groups[*].name : v == local.name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
