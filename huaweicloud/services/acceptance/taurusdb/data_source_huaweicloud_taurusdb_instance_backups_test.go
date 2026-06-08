package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBInstanceBackups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_instance_backups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBInstanceBackups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size_unit"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.backup_type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.backup_level"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.backup_method"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.time_zone"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("order_field_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBInstanceBackups_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_instance_backups" "test" {
  instance_id = "%[1]s"
}

locals {
  backup_name = data.huaweicloud_taurusdb_instance_backups.test.backups[0].name
}

data "huaweicloud_taurusdb_instance_backups" "filter" {
  instance_id    = "%[1]s"
  filter_field   = "name"
  filter_content = local.backup_name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_instance_backups.filter.backups) > 0 && alltrue([
    for v in data.huaweicloud_taurusdb_instance_backups.filter.backups[*].name : v == local.backup_name
  ])
}

data "huaweicloud_taurusdb_instance_backups" "order" {
  instance_id = "%[1]s"
  order_field = "beginTime"
  order_rule  = "asc"
}

locals {
  origin_first_backup_created   = data.huaweicloud_taurusdb_instance_backups.test.backups[0].created
  latest_backups_index          = length(data.huaweicloud_taurusdb_instance_backups.order.backups) - 1
  ordered_latest_backup_created = data.huaweicloud_taurusdb_instance_backups.order.backups[local.latest_backups_index].created
}

output "order_field_is_useful" {
  value = local.origin_first_backup_created == local.ordered_latest_backup_created
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID)
}
