package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccProtectedInstanceDeleteNIC_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSDRSDeleteNic(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProtectedInstanceDeleteNIC_basic(),
			},
		},
	})
}

func testProtectedInstanceDeleteNIC_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sdrs_protected_instance_delete_nic" "test" {
  protected_instance_id = "%[1]s"
  nic_id                = "%[2]s"
}
`, acceptance.HW_SDRS_PROTECTION_INSTANCE_ID, acceptance.HW_SDRS_NIC_ID)
}
