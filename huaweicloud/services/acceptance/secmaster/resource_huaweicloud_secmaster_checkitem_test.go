package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getResourceCheckitemFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetCheckitemInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceCheckitem_basic(t *testing.T) {
	var (
		rName         = "huaweicloud_secmaster_checkitem.test"
		checkitemName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceCheckitemFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCheckitem_basic(checkitemName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", checkitemName),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "level", "medium"),
					resource.TestCheckResourceAttr(rName, "cloud_server", "IAM"),
					resource.TestCheckResourceAttr(rName, "method", "0"),
					resource.TestCheckResourceAttr(rName, "source", "3"),
				),
			},
			{
				Config: testAccResourceCheckitem_update(checkitemName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", checkitemName)),
					resource.TestCheckResourceAttr(rName, "description", "updated description"),
					resource.TestCheckResourceAttr(rName, "level", "high"),
					resource.TestCheckResourceAttr(rName, "cloud_server", "ECS"),
					resource.TestCheckResourceAttr(rName, "method", "1"),
					resource.TestCheckResourceAttr(rName, "source", "3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCheckitemImportStateFunc(rName),
			},
		},
	})
}

func testAccResourceCheckitem_basic(checkitemName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_checkitem" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  description  = "test description"
  level        = "medium"
  cloud_server = "IAM"
  method       = 0
  source       = 3
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, checkitemName)
}

func testAccResourceCheckitem_update(checkitemName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_checkitem" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_update"
  description  = "updated description"
  level        = "high"
  cloud_server = "ECS"
  method       = 1
  source       = 3
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, checkitemName)
}

func testAccCheckitemImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, name string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		name = rs.Primary.ID

		if workspaceId == "" || name == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<name>', but got '%s/%s'",
				workspaceId, name)
		}

		return fmt.Sprintf("%s/%s", workspaceId, name), nil
	}
}
