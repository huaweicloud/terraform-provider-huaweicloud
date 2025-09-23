package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcVirtualInterfaceSwitchoverRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dc_virtual_interface_switchover_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcVirtualInterfaceSwitchoverRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.operation"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.operate_status"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "switchover_test_records.0.end_time"),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcVirtualInterfaceSwitchoverRecords_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_virtual_interface_switchover_records" "test" {
  depends_on = [huaweicloud_dc_virtual_interface_switchover.test]
}

locals {
  resource_id    = huaweicloud_dc_virtual_interface.test.id
  filter_records = data.huaweicloud_dc_virtual_interface_switchover_records.resource_id_filter.switchover_test_records
}
data "huaweicloud_dc_virtual_interface_switchover_records" "resource_id_filter" {
  depends_on = [huaweicloud_dc_virtual_interface_switchover.test]

  resource_id = [huaweicloud_dc_virtual_interface.test.id]
}
output "resource_id_filter_is_useful" {
  value = length(local.filter_records) > 0 && alltrue([for v in local.filter_records[*].resource_id : v == local.resource_id]
  )
}
`, testVirtualInterfaceSwitchover_basic(name))
}
