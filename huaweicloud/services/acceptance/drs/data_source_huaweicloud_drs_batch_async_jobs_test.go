package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBatchAsyncJobs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_batch_async_jobs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBatchAsyncJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "jobs.#"),
				),
			},
		},
	})
}

// This resource has no response data, and the filtering parameters cannot be validated in the test case.
func testDataSourceBatchAsyncJobs_basic() string {
	return `
data "huaweicloud_drs_batch_async_jobs" "test" {
  status   = "ASYNC_JOB_VALIDATING"
  sort_key = "create_time"
  sort_dir = "desc"
}
`
}
