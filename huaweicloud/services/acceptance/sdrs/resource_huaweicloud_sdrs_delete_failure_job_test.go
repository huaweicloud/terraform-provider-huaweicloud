package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeleteFailureJob_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSDRSFailureJob(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDeleteFailureJob_basic(),
			},
		},
	})
}

func testDeleteFailureJob_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sdrs_delete_failure_job" "test" {
  failure_job_id = "%[1]s"
}
`, acceptance.HW_SDRS_FAILURE_JOB_ID)
}
