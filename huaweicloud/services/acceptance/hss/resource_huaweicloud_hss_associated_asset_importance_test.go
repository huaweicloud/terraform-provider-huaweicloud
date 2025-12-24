package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAssociatedAssetImportance_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAssociatedAssetImportance_basic(),
			},
		},
	})
}

func testAssociatedAssetImportance_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_associated_asset_importance" "test" {
  asset_value           = "important"
  host_id_list          = ["%s"]
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
