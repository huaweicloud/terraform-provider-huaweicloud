package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVirtualInterfaceAccepter_accepted(t *testing.T) {
	// The resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// The environment variable `HW_DC_VIRTUAL_INTERFACE_ID` refers to the virtual interface ID created by other
			// tenants and under the same region.
			acceptance.TestAccPreCheckDCVirtualInterfaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualInterfaceAccepter_action("ACCEPTED"),
			},
		},
	})
}

func TestAccVirtualInterfaceAccepter_rejected(t *testing.T) {
	// The resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// The environment variable `HW_DC_VIRTUAL_INTERFACE_ID` refers to the virtual interface ID created by other
			// tenants and under the same region.
			acceptance.TestAccPreCheckDCVirtualInterfaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualInterfaceAccepter_action("REJECTED"),
			},
		},
	})
}

func testAccVirtualInterfaceAccepter_action(action string) string {
	return fmt.Sprintf(`

resource "huaweicloud_dc_virtual_interface_accepter" "test" {
  virtual_interface_id = "%[1]s"
  action               = "%[2]s"
}
`, acceptance.HW_DC_VIRTUAL_INTERFACE_ID, action)
}
