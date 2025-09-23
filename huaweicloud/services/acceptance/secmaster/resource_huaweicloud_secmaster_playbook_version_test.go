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

func getPlaybookVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	workspaceID := state.Primary.Attributes["workspace_id"]
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetPlaybookVersion(client, workspaceID, state.Primary.ID)
}

func TestAccPlaybookVersion_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_playbook_version.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPlaybookVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "playbook_id",
						"huaweicloud_secmaster_playbook.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "creator_id"),
					resource.TestCheckResourceAttrSet(rName, "enabled"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testPlaybookVersion_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "playbook_id",
						"huaweicloud_secmaster_playbook.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(rName, "dataobject_create", "true"),
					resource.TestCheckResourceAttr(rName, "dataobject_update", "true"),
					resource.TestCheckResourceAttr(rName, "dataobject_delete", "true"),
					resource.TestCheckResourceAttr(rName, "rule_enable", "true"),
					resource.TestCheckResourceAttr(rName, "trigger_type", "EVENT"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookVersionImportState(rName),
			},
		},
	})
}

// The playbook version can be updated after matching a workflow (playbook action).
func testPlaybookVersion_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id  = "%[2]s"
  business_code = "Policy"
}

resource "huaweicloud_secmaster_playbook_version" "test" {
  workspace_id      = "%[2]s"
  playbook_id       = huaweicloud_secmaster_playbook.test.id
  dataclass_id      = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  description       = "created by terraform"
  dataobject_create = true
  trigger_type      = "EVENT"
}

data "huaweicloud_secmaster_workflows" "test" {
  workspace_id  = "%[2]s"
  data_class_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
}
  
resource "huaweicloud_secmaster_playbook_action" "test" {
  workspace_id = "%[2]s"
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  action_id    = data.huaweicloud_secmaster_workflows.test.workflows[0].id
  name         = data.huaweicloud_secmaster_workflows.test.workflows[0].name
  description  = "created by terraform"
}
`, testPlaybook_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testPlaybookVersion_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_data_classes" "test" {
  workspace_id = "%[2]s"
}

resource "huaweicloud_secmaster_playbook_version" "test" {
  workspace_id      = "%[2]s"
  playbook_id       = huaweicloud_secmaster_playbook.test.id
  dataclass_id      = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
  description       = "updated by terraform"
  dataobject_create = true
  dataobject_update = true
  dataobject_delete = true
  rule_enable       = true
  trigger_type      = "EVENT"
}

data "huaweicloud_secmaster_workflows" "test" {
  workspace_id  = "%[2]s"
  data_class_id = data.huaweicloud_secmaster_data_classes.test.data_classes[0].id
}

resource "huaweicloud_secmaster_playbook_action" "test" {
  workspace_id = "%[2]s"
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  action_id    = data.huaweicloud_secmaster_workflows.test.workflows[0].id
  name         = data.huaweicloud_secmaster_workflows.test.workflows[0].name
  description  = "created by terraform"
}
`, testPlaybook_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testPlaybookVersionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["playbook_id"], rs.Primary.ID), nil
	}
}
