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

func getResourceWorkflowFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetWorkflowInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceWorkflow_basic(t *testing.T) {
	var (
		rName      = "huaweicloud_secmaster_workflow.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceWorkflowFunc,
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
				Config: testAccWorkflow_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "dataclass_id", "data.huaweicloud_secmaster_data_classes.test", "data_classes.0.id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "engine_type", "public_engine"),
					resource.TestCheckResourceAttr(rName, "aop_type", "NORMAL"),
					resource.TestCheckResourceAttr(rName, "labels", "IP"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttrSet(rName, "dataclass_name"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
				),
			},
			{
				Config: testAccWorkflow_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "labels", "ACCOUNT"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWorkflowImportStateFunc(rName),
			},
		},
	})
}

func testAccWorkflow_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = "%[1]s"
}

resource "huaweicloud_secmaster_workflow" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  dataclass_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  engine_type  = "public_engine"
  aop_type     = "NORMAL"
  labels       = "IP"
  description  = "terraform test"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testAccWorkflow_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = "%[1]s"
}

resource "huaweicloud_secmaster_workflow" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  dataclass_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  engine_type  = "public_engine"
  aop_type     = "NORMAL"
  labels       = "ACCOUNT"
  description  = ""
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, name)
}

func testAccWorkflowImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, workflowId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		workflowId = rs.Primary.ID

		if workspaceId == "" || workflowId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, workflowId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, workflowId), nil
	}
}
