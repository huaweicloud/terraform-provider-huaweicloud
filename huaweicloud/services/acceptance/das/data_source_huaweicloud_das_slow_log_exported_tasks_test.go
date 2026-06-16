package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSlowLogExportedTasks_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_slow_log_exported_tasks.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByExportType   = "data.huaweicloud_das_slow_log_exported_tasks.filter_by_export_type"
		dcFilterByExportType = acceptance.InitDataSourceCheck(filterByExportType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSlowLogExportedTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tasks.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.created_time"),

					// filter by export_type
					dcFilterByExportType.CheckResourceExists(),
					resource.TestCheckOutput("is_export_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSlowLogExportedTasks_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccDataSlowLogExportedTasks_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_slow_log_exported_tasks" "all" {
  instance_id = local.instance_ids[0]
}

# Filter by export_type
locals {
  export_type = "slowsqldetails"
}

data "huaweicloud_das_slow_log_exported_tasks" "filter_by_export_type" {
  instance_id = local.instance_ids[0]
  export_type = local.export_type
}

output "is_export_type_filter_useful" {
  value = length(data.huaweicloud_das_slow_log_exported_tasks.filter_by_export_type.tasks) > 0
}
`, testAccDataSlowLogExportedTasks_base())
}
