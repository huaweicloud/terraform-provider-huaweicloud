package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCloudLogResource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCloudLogResource_basic(),
			},
		},
	})
}

func testCloudLogResource_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_cloud_log_resource" "test" {
  workspace_id = "%[1]s"
  domain_id    = "%[2]s"

  resources {
    enable    = "active"
    region_id = "%[3]s"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_DOMAIN_ID, acceptance.HW_REGION_NAME)
}
