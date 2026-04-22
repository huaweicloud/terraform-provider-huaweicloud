package dcs

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceBackgroundTasks_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_background_tasks.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackgroundTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "tasks.#"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.#"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.old_capacity"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.new_capacity"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.enable_public_ip"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.public_ip_id"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.public_ip_address"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.enable_ssl"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.old_cache_mode"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.new_cache_mode"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.old_resource_spec_code"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.new_resource_spec_code"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.old_replica_num"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.new_replica_num"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.old_cache_type"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.new_cache_type"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.replica_ip"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.replica_az"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.group_name"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.old_port"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.new_port"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.is_only_adjust_charging"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.account_name"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.source_ip"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.target_ip"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.node_name"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.rename_commands.#"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.details.0.updated_config_length"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.user_name"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.params"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.error_code"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.enable_show"),
					resource.TestCheckResourceAttrSet(rName, "tasks.0.job_id"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceBackgroundTasks_basic() string {
	startTime := time.Now().UTC().Add(-1 * time.Hour).Format("20060102150405")
	endTime := time.Now().UTC().Add(1 * time.Hour).Format("20060102150405")
	return fmt.Sprintf(`
data "huaweicloud_dcs_background_tasks" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_dcs_background_tasks" "time_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_dcs_background_tasks.time_filter.tasks) > 0
}
`, acceptance.HW_DCS_INSTANCE_ID, startTime, endTime)
}
