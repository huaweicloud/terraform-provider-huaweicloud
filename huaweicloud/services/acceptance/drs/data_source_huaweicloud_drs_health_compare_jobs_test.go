package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDrsHealthCompareJobs_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_drs_health_compare_jobs.test"
		dc         = acceptance.InitDataSourceCheck(datasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDrsHealthCompareJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "compare_jobs.#"),
				),
			},
		},
	})
}

func testAccDatasourceDrsHealthCompareJobs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_health_compare_jobs" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
