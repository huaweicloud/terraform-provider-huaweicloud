package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSRecyclePolicy_basic(t *testing.T) {
	resourceName := "huaweicloud_dds_recycle_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSRecyclePolicy_basic(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "retention_period_in_days", "1"),
				),
			},
			{
				Config: testAccDDSRecyclePolicy_basic(6),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "retention_period_in_days", "6"),
				),
			},
		},
	})
}

func testAccDDSRecyclePolicy_basic(day int) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_recycle_policy" "test" {
  retention_period_in_days = %d
}`, day)
}
