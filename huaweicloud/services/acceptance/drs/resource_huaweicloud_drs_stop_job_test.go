package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccStopJob_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testStopJob_basic(),
			},
		},
	})
}

func testStopJob_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_stop_job" "test" {
  job_id        = "%s"
  is_force_stop = true
}
`, acceptance.HW_DRS_JOB_ID)
}
