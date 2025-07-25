package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBackupPolicy_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_backup_policy.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBackupPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_type"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.day_backups"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.max_backups"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.month_backups"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.retention_duration_days"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.timezone"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.week_backups"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_definition.0.year_backups"),
					resource.TestCheckResourceAttrSet(dataSource, "trigger.#"),
				),
			},
		},
	})
}

func testAccDataSourceBackupPolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_policy" "test" {
  name        = "%s"
  type        = "backup"
  time_period = 20

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}

data "huaweicloud_hss_backup_policy" "test" {
  policy_id             = huaweicloud_cbr_policy.test.id
  enterprise_project_id = "all_granted_eps"
}
`, acceptance.RandomAccResourceName())
}
