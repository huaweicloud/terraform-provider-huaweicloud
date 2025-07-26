package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance, please ensure the Workspace service is enabled.
// lintignore:AT001
func TestAccAppServiceAction_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAction_basic,
			},
		},
	})
}

const testAccAppServiceAction_basic = `
resource "huaweicloud_workspace_app_service_action" "test" {
  service_status = "active"
}
`
