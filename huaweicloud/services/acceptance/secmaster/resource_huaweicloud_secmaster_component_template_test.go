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

func getResourceComponentTemplateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetComponentTemplateInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceComponentTemplate_basic(t *testing.T) {
	var (
		rName      = "huaweicloud_secmaster_component_template.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceComponentTemplateFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterComponentId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponentTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "component_id", acceptance.HW_SECMASTER_COMPONENT_ID),
					resource.TestCheckResourceAttr(rName, "template_name", name),
					resource.TestCheckResourceAttr(rName, "task_config", "{\"action_id\":\"test_action\"}"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testAccComponentTemplate_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "template_name", updateName),
					resource.TestCheckResourceAttr(rName, "task_config", "{\"action_id\":\"update_action\"}"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccComponentTemplateImportStateFunc(rName),
			},
		},
	})
}

func testAccComponentTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_component_template" "test" {
  workspace_id  = "%[1]s"
  component_id  = "%[2]s"
  template_name = "%[3]s"
  task_config   = "{\"action_id\":\"test_action\"}"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID, name)
}

func testAccComponentTemplate_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_component_template" "test" {
  workspace_id  = "%[1]s"
  component_id  = "%[2]s"
  template_name = "%[3]s"
  task_config   = "{\"action_id\":\"update_action\"}"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_COMPONENT_ID, name)
}

func testAccComponentTemplateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, templateId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		templateId = rs.Primary.ID

		if workspaceId == "" || templateId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, templateId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, templateId), nil
	}
}
