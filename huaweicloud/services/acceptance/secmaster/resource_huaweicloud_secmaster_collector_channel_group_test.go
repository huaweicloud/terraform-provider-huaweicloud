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

func getResourceCollectorChannelGroupFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetCollectorChannelGroupByName(client, state.Primary.Attributes["workspace_id"], state.Primary.Attributes["name"])
}

func TestAccResourceCollectorChannelGroup_basic(t *testing.T) {
	var (
		rName      = "huaweicloud_secmaster_collector_channel_group.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceCollectorChannelGroupFunc,
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
				Config: testAccResourceCollectorChannelGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
				),
			},
			{
				Config: testAccResourceCollectorChannelGroup_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCollectorChannelGroupImportStateFunc(rName),
			},
		},
	})
}

func testAccResourceCollectorChannelGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collector_channel_group" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testAccResourceCollectorChannelGroup_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collector_channel_group" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testAccCollectorChannelGroupImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, groupName string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		groupName = rs.Primary.Attributes["name"]

		if workspaceId == "" || groupName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<name>', but got '%s/%s'",
				workspaceId, groupName)
		}

		return fmt.Sprintf("%s/%s", workspaceId, groupName), nil
	}
}
