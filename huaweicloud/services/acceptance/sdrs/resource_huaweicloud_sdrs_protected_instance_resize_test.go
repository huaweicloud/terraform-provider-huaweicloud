package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccProtectedInstanceResize_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSDRSInstanceID(t)
			acceptance.TestAccPreCheckSDRSInstanceResizeFlavor(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProtectedInstanceResize_basic(),
			},
		},
	})
}

func testProtectedInstanceResize_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sdrs_protected_instance_resize" "test" {
  protected_instance_id = "%[1]s"
  flavor_ref            = "%[2]s"
}
`, acceptance.HW_SDRS_PROTECTION_INSTANCE_ID, acceptance.HW_SDRS_RESIZE_FLAVOR_ID)
}
