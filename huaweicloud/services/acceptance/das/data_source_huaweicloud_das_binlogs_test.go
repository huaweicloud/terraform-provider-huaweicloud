package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataBinlogs_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_binlogs.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBinlogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "binlogs.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "binlogs.0.file_name"),
					resource.TestCheckResourceAttrSet(all, "binlogs.0.file_size"),
				),
			},
		},
	})
}

func TestAccDataBinlogs_attributes(t *testing.T) {
	var (
		dcName = "data.huaweicloud_das_binlogs.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBinlogs_attributes(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "binlogs.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.file_name"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.file_size"),
					resource.TestMatchResourceAttr(dcName, "binlogs.0.task_info.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.binlog_type"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.connection_id"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.file_name"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.project_id"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.project_name"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.status"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.user_id"),
					resource.TestCheckResourceAttrSet(dcName, "binlogs.0.task_info.0.user_name"),
					resource.TestMatchResourceAttr(dcName, "binlogs.0.task_info.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dcName, "binlogs.0.task_info.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataBinlogs_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}

data "huaweicloud_das_database_users" "test" {
  instance_id = local.instance_ids[0]
}

locals {
  user_id = try(data.huaweicloud_das_database_users.test.users.0.id, "")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccDataBinlogs_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_binlogs" "all" {
  user_id     = local.user_id
  binlog_type = "latest"
}
`, testAccDataBinlogs_base())
}

func testAccDataBinlogs_attributes() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_binlogs" "test" {
  user_id     = local.user_id
  binlog_type = "latest"
}

# If no binlog is parsed, the 'task_info' field will be empty.
# So we need to manully create a binlog parse task.
resource "huaweicloud_das_binlog_parse_task" "test" {
  user_id     = local.user_id
  binlog_type = "latest"
  file_name   = try(data.huaweicloud_das_binlogs.test.binlogs.0.file_name, "")
}
`, testAccDataBinlogs_base())
}
