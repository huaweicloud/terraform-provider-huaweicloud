package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppAction_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_action.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is used to manage the tenant profile settings. Deleting this resource will not
		// clear the corresponding configuration, but will only remove the resource information from the tfstate file.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppAction_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_restrict_rule_switch", "true"),
				),
			},
		},
	})
}

func testAccAppAction_basic() string {
	return `
resource "huaweicloud_workspace_app_action" "test" {
  app_restrict_rule_switch = true
}
`
}
