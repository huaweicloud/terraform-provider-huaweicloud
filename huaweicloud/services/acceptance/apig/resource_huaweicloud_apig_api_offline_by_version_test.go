package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigApiOfflineByVersion_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_apig_api_offline_by_version.test"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApigApiOfflineByVersion_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccApigApiOfflineByVersion_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_api_offline_by_version" "test" {
  instance_id = "%s"
  version_id  = "%s"
}
`, acceptance.HW_APIG_INSTANCE_ID, acceptance.HW_APIG_API_VERSION_ID)
}
