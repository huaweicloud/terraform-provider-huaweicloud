package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataspace_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDataspace_basic(name),
			},
		},
	})
}

func testAccDataspace_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = "%[1]s"
  dataspace_name = "%[2]s"
  description    = "test description"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}
