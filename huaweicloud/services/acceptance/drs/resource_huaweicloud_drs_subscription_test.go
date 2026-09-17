package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrsSubscription_basic(t *testing.T) {
	resourceName := "huaweicloud_drs_subscription.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsSubscription_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "DRS-123"),
				),
			},
		},
	})
}

func testAccDrsSubscription_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_drs_subscription" "test" {
  name                  = "DRS-123"
  enterprise_project_id = "0"

  source_endpoint_info {
    id = "%[1]s"
  }
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
