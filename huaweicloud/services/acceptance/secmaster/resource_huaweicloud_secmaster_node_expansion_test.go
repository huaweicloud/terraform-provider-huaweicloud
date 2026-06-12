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

func getResourceNodeExpansionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetNodeExpansionByNodeId(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceNodeExpansion_basic(t *testing.T) {
	var (
		rName = "huaweicloud_secmaster_node_expansion.test"

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceNodeExpansionFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodeExpansion_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "node_id", acceptance.HW_SECMASTER_NODE_ID),
					resource.TestCheckResourceAttr(rName, "custom_label", "test_label"),
					resource.TestCheckResourceAttr(rName, "data_center", "test_center"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "maintainer", "admin"),
					resource.TestCheckResourceAttr(rName, "network_plane", "business"),
				),
			},
			{
				Config: testAccNodeExpansion_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "custom_label", "update_label"),
					resource.TestCheckResourceAttr(rName, "data_center", "update_center"),
					resource.TestCheckResourceAttr(rName, "description", "update description"),
					resource.TestCheckResourceAttr(rName, "maintainer", "operator"),
					resource.TestCheckResourceAttr(rName, "network_plane", "management"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNodeExpansionImportStateFunc(rName),
			},
		},
	})
}

func testAccNodeExpansion_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_node_expansion" "test" {
  workspace_id   = "%[1]s"
  node_id        = "%[2]s"
  custom_label   = "test_label"
  data_center    = "test_center"
  description    = "test description"
  maintainer     = "admin"
  network_plane  = "business"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_NODE_ID)
}

func testAccNodeExpansion_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_node_expansion" "test" {
  workspace_id   = "%[1]s"
  node_id        = "%[2]s"
  custom_label   = "update_label"
  data_center    = "update_center"
  description    = "update description"
  maintainer     = "operator"
  network_plane  = "management"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_NODE_ID)
}

func testAccNodeExpansionImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, nodeId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		nodeId = rs.Primary.ID

		if workspaceId == "" || nodeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, nodeId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, nodeId), nil
	}
}
