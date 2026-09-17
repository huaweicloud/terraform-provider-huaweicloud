package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsSubscriptionReset_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_subscription_reset.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsSubscriptionReset_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_DRS_JOB_ID),
				),
			},
		},
	})
}

func testAccDrsSubscriptionReset_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_subscription_reset" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
