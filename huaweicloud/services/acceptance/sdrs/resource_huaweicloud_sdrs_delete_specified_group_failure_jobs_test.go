package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeleteSpecifiedGroupFailureJobs_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSDRSProtectionGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDeleteSpecifiedGroupFailureJobs_basic(),
			},
		},
	})
}

func testDeleteSpecifiedGroupFailureJobs_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sdrs_delete_specified_group_failure_jobs" "test" {
  server_group_id = "%[1]s"
}
`, acceptance.HW_SDRS_PROTECTION_GROUP_ID)
}
