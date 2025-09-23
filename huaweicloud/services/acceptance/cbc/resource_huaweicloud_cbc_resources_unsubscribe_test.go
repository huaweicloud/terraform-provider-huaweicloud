package cbc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourcesUnsubscribe_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time resource and has not destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCbcResourcesUnsubscribe(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcesUnsubscribe_basic(),
			},
		},
	})
}

func testAccResourcesUnsubscribe_basic() string {
	ramdomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
resource "huaweicloud_cbc_resources_unsubscribe" "test" {
  resource_ids = [
    "%[1]s",
    "%[2]s" # Providing a non-existent resource.
  ]

  enable_force_new = true
}
`, acceptance.HW_CBC_UNSUBSCRIBE_RESOURCE_ID, ramdomId)
}
