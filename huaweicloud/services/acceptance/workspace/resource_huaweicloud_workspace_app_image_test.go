package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppImage_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_image.test"
		name         = acceptance.RandomAccResourceName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
			acceptance.TestAccPrecheckWorkspaceUserNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceAppImage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Workspace APP genereted image"),
				),
			},
		},
	})
}

func testResourceAppImage_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_image_server" "test" {
  name                    = "%[1]s"
  flavor_id               = "%[2]s"
  vpc_id                  = "%[3]s"
  subnet_id               = "%[4]s"
  image_id                = "%[5]s"
  image_type              = "gold"
  image_source_product_id = "%[6]s"
  spec_code               = "%[7]s"

  # Currently only one user can be set.
  authorize_accounts {
    account = split(",", "%[8]s")[0]
    type    = "USER"
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  is_vdi                         = true
  enterprise_project_id          = "%[9]s"
  is_delete_associated_resources = true
}

resource "huaweicloud_workspace_app_image" "test" {
  server_id             = huaweicloud_workspace_app_image_server.test.id
  name                  = "%[1]s"
  description           = "Workspace APP genereted image"
  enterprise_project_id = "%[9]s"
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_AD_VPC_ID,
		acceptance.HW_WORKSPACE_AD_NETWORK_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE,
		acceptance.HW_WORKSPACE_USER_NAMES,
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
	)
}
