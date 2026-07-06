package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscScanJobs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dsc_scan_jobs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byContent   = "data.huaweicloud_dsc_scan_jobs.content_filter"
		dcByContent = acceptance.InitDataSourceCheck(byContent)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscScanJobs_conf,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.cycle"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.0.create_time"),

					dcByContent.CheckResourceExists(),
					resource.TestCheckOutput("is_content_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceDscScanJobs_conf = `
data "huaweicloud_dsc_scan_jobs" "test" {
}

# Filter by content
locals {
  content = data.huaweicloud_dsc_scan_jobs.test.jobs[0].name
}

data "huaweicloud_dsc_scan_jobs" "content_filter" {
  content = local.content
}

locals {
  content_filter_result = [
    for v in data.huaweicloud_dsc_scan_jobs.content_filter.jobs[*].name : v == local.content
  ]
}

output "is_content_filter_useful" {
  value = alltrue(local.content_filter_result) && length(local.content_filter_result) > 0
}
`
