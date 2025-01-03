package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPlaybookEnableResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	workspaceID := state.Primary.Attributes["workspace_id"]
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	playbook, err := secmaster.GetPlaybook(client, workspaceID, state.Primary.ID)
	if err != nil {
		return nil, err
	}

	enabled := utils.PathSearch("enabled", playbook, false).(bool)
	if !enabled {
		return nil, golangsdk.ErrDefault404{}
	}

	return playbook, err
}

func TestAccPlaybookEnable_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_playbook_enable.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookEnableResourceFunc,
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
				Config: testPlaybookEnable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "playbook_id", "huaweicloud_secmaster_playbook.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "playbook_name", "huaweicloud_secmaster_playbook.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "active_version_id", "huaweicloud_secmaster_playbook_approval.test", "id"),
				),
			},
		},
	})
}

func testPlaybookEnable_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_secmaster_playbook_version_action" "submit_version" {
  workspace_id = "%[2]s"
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  status       = "APPROVING"

  depends_on = [huaweicloud_secmaster_playbook_action.test]

  lifecycle {
    ignore_changes = [
      status, enabled,
    ]
  }
}

resource "huaweicloud_secmaster_playbook_approval" "test" {
  workspace_id = "%[2]s"
  version_id   = huaweicloud_secmaster_playbook_version.test.id
  result       = "PASS"
  content      = "ok"

  depends_on = [huaweicloud_secmaster_playbook_version_action.submit_version]
}

resource "huaweicloud_secmaster_playbook_enable" "test" {
  workspace_id      = "%[2]s"
  playbook_id       = huaweicloud_secmaster_playbook.test.id
  playbook_name     = huaweicloud_secmaster_playbook.test.name
  active_version_id = huaweicloud_secmaster_playbook_approval.test.id
}
`, testPlaybookVersion_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}
