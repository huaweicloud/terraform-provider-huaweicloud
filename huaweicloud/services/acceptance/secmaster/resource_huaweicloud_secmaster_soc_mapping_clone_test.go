package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSocMappingClone_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterMappingId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSocMappingClone_basic(name),
			},
		},
	})
}

func testAccSocMappingClone_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_soc_mapping_clone" "test" {
  workspace_id = "%[1]s"
  mapping_id   = "%[2]s"
  name         = "%[3]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_MAPPING_ID, name)
}
