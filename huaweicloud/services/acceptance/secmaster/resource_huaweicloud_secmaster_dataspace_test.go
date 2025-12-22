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

func getResourceDataspaceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetDataspaceInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccDataspace_basic(t *testing.T) {
	var (
		rName = "huaweicloud_secmaster_dataspace.test"
		name  = acceptance.RandomAccResourceNameWithDash()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceDataspaceFunc,
		)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataspace_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "dataspace_name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttrSet(rName, "dataspace_type"),
					resource.TestCheckResourceAttrSet(rName, "create_by"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				Config: testAccDataspace_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", "test terraform"),
					resource.TestCheckResourceAttrSet(rName, "update_by"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDataspaceImportStateFunc(rName),
			},
		},
	})
}

func testAccDataspace_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = "%[1]s"
  dataspace_name = "%[2]s"
  description    = "test description"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testAccDataspace_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_dataspace" "test" {
  workspace_id   = "%[1]s"
  dataspace_name = "%[2]s"
  description    = "test terraform"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testAccDataspaceImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, dataspaceId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		dataspaceId = rs.Primary.ID

		if workspaceId == "" || dataspaceId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, dataspaceId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, dataspaceId), nil
	}
}
