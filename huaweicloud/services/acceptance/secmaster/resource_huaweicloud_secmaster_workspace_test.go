package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkspace_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspace(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testWorkspace_basic(name),
			},
		},
	})
}

func testWorkspace_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_workspace" "test" {
  name         = "%[1]s"
  project_name = "%[2]s"
}
`, name, acceptance.HW_REGION_NAME)
}
