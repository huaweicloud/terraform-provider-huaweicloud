package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClassTypeLayout_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterDataClassID(t)
			acceptance.TestAccPreCheckSecMasterTypeID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testDataClassTypeLayout_basic(),
			},
		},
	})
}

func testDataClassTypeLayout_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_data_class_type_layout" "test" {
  workspace_id = "%[1]s"
  dataclass_id = "%[2]s"
  type_id      = "%[3]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_DATACLASS_ID, acceptance.HW_SECMASTER_TYPE_ID)
}
