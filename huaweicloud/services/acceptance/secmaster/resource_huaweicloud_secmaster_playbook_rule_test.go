package secmaster

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPlaybookRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPlaybookRule: Query the SecMaster playbook detail
	var (
		getPlaybookRuleHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/rules/{id}"
		getPlaybookRuleProduct = "secmaster"
	)
	getPlaybookRuleClient, err := cfg.NewServiceClient(getPlaybookRuleProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookRulePath := getPlaybookRuleClient.Endpoint + getPlaybookRuleHttpUrl
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{project_id}", getPlaybookRuleClient.ProjectID)
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{workspace_id}", state.Primary.Attributes["workspace_id"])
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{version_id}", state.Primary.Attributes["version_id"])
	getPlaybookRulePath = strings.ReplaceAll(getPlaybookRulePath, "{id}", state.Primary.ID)

	getPlaybookRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookRuleResp, err := getPlaybookRuleClient.Request("GET", getPlaybookRulePath, &getPlaybookRuleOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getPlaybookRuleResp)
}

func TestAccPlaybookRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_secmaster_playbook_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPlaybookRuleResourceFunc,
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
				Config: testPlaybookRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "version_id",
						"huaweicloud_secmaster_playbook_version.test", "id"),
					resource.TestCheckResourceAttr(rName, "expression_type", "common"),
					resource.TestCheckResourceAttr(rName, "conditions.0.name", "condition_0"),
					resource.TestCheckResourceAttr(rName, "conditions.0.detail", "123"),
					resource.TestCheckResourceAttr(rName, "logics.0", "condition_0"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testPlaybookRule_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttrPair(rName, "version_id",
						"huaweicloud_secmaster_playbook_version.test", "id"),
					resource.TestCheckResourceAttr(rName, "expression_type", "common"),
					resource.TestCheckResourceAttr(rName, "conditions.1.name", "condition_1"),
					resource.TestCheckResourceAttr(rName, "conditions.1.detail", "456"),
					resource.TestCheckResourceAttr(rName, "logics.1", "condition_1"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlaybookRuleImportState(rName),
			},
		},
	})
}

func testPlaybookRule_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_secmaster_playbook_rule" "test" {
  workspace_id    = "%s"
  version_id      = huaweicloud_secmaster_playbook_version.test.id
  expression_type = "common"

  conditions {
    name   = "condition_0"
    detail = "123"
    data   = ["waf.alarm.level", ">", "3"]
  }

  logics = ["condition_0"]
}
`, testPlaybookVersion_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testPlaybookRule_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_secmaster_playbook_rule" "test" {
  workspace_id    = "%s"
  version_id      = huaweicloud_secmaster_playbook_version.test.id
  expression_type = "common"

  conditions {
    name   = "condition_0"
    detail = "123"
    data   = ["waf.alarm.level", ">", "3"]
  }

  conditions {
    name   = "condition_1"
    detail = "456"
    data   = ["waf.alarm.level", "=", "4"]
  }

  logics = ["condition_0","condition_1"]
}
`, testPlaybookVersion_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}

func testPlaybookRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["workspace_id"] == "" {
			return "", fmt.Errorf("attribute (workspace_id) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["workspace_id"], rs.Primary.Attributes["version_id"], rs.Primary.ID), nil
	}
}
