package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInstanceProcesses_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_instance_processes.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByUser   = "data.huaweicloud_das_instance_processes.filter_by_user"
		dcFilterByUser = acceptance.InitDataSourceCheck(filterByUser)

		filterByDatabase   = "data.huaweicloud_das_instance_processes.filter_by_database"
		dcFilterByDatabase = acceptance.InitDataSourceCheck(filterByDatabase)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceProcesses_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "processes.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "processes.0.id"),
					resource.TestCheckResourceAttrSet(all, "processes.0.db_user_name"),
					resource.TestCheckResourceAttrSet(all, "processes.0.host"),
					resource.TestCheckResourceAttrSet(all, "processes.0.db_name"),
					resource.TestCheckResourceAttrSet(all, "processes.0.command"),
					resource.TestCheckResourceAttrSet(all, "processes.0.time"),
					resource.TestCheckResourceAttrSet(all, "processes.0.trx_executed_time"),

					// filter by user
					dcFilterByUser.CheckResourceExists(),
					resource.TestCheckOutput("is_user_filter_useful", "true"),

					// filter by database
					dcFilterByDatabase.CheckResourceExists(),
					resource.TestCheckOutput("is_database_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccInstanceProcesses_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}

data "huaweicloud_das_database_users" "test" {
  instance_id = local.instance_ids[0]
}

locals {
  db_user_id = try(data.huaweicloud_das_database_users.test.users[0].id, "")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccInstanceProcesses_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_instance_processes" "all" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
}

# Filter by user
locals {
  db_user_name = data.huaweicloud_das_instance_processes.all.processes[0].db_user_name
}

data "huaweicloud_das_instance_processes" "filter_by_user" {
  instance_id  = local.instance_ids[0]
  db_user_id   = local.db_user_id
  db_user_name = local.db_user_name
}

locals {
  user_filter_result = [
    for v in data.huaweicloud_das_instance_processes.filter_by_user.processes : v.db_user_name == local.db_user_name
  ]
}

output "is_user_filter_useful" {
  value = length(local.user_filter_result) > 0 && alltrue(local.user_filter_result)
}

# Filter by database
locals {
  db_name = data.huaweicloud_das_instance_processes.all.processes[0].db_name
}

data "huaweicloud_das_instance_processes" "filter_by_database" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
  db_name     = local.db_name
}

locals {
  database_filter_result = [
    for v in data.huaweicloud_das_instance_processes.filter_by_database.processes : v.db_name == local.db_name
  ]
}

output "is_database_filter_useful" {
  value = length(local.database_filter_result) > 0 && alltrue(local.database_filter_result)
}
`, testAccInstanceProcesses_base())
}
