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

func getResourcePipeConsumptionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetPipeConsumptionByName(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourcePipeConsumption_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_secmaster_pipe_consumption.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourcePipeConsumptionFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterPipeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePipeConsumption_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "pipe_id", acceptance.HW_SECMASTER_PIPE_ID),
					resource.TestCheckResourceAttrSet(rName, "access_point"),
					resource.TestCheckResourceAttrSet(rName, "pipe_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccPipeConsumptionImportStateFunc(rName),
			},
		},
	})
}

func testAccResourcePipeConsumption_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_pipe_consumption" "test" {
  workspace_id = "%[1]s"
  pipe_id      = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_PIPE_ID)
}

func testAccPipeConsumptionImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, pipeId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		pipeId = rs.Primary.Attributes["pipe_id"]

		if workspaceId == "" || pipeId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<pipe_id>', but got '%s/%s'",
				workspaceId, pipeId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, pipeId), nil
	}
}
