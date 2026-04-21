package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getSecurityDataRecognitionRuleGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	return dataarts.GetSecurityDataRecognitionRuleGroupById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccSecurityDataRecognitionRuleGroup_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_security_data_recognition_rule_group.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getSecurityDataRecognitionRuleGroupResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsSecurityDataCategoryIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccSecurityDataRecognitionRuleGroup_nonExistentWorkspaceAndRule(name),
				ExpectError: regexp.MustCompile("error creating DataArts Security data recognition rule group"),
			},
			{
				Config: testAccSecurityDataRecognitionRuleGroup_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "rule_ids.0",
						"huaweicloud_dataarts_security_data_recognition_rule.test.0", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccSecurityDataRecognitionRuleGroup_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "rule_ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "rule_ids.0",
						"huaweicloud_dataarts_security_data_recognition_rule.test.1", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSecurityDataRecognitionRuleGroupImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccSecurityDataRecognitionRuleGroupImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		id := rs.Primary.ID
		if workspaceID == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceID, id)
		}
		return fmt.Sprintf("%s/%s", workspaceID, id), nil
	}
}

func testAccSecurityDataRecognitionRuleGroup_nonExistentWorkspaceAndRule(name string) string {
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_recognition_rule_group" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  rule_ids     = ["%[1]s"]
}
`, randomUUID, name)
}

func testAccSecurityDataRecognitionRuleGroup_basic_base(name string) string {
	return fmt.Sprintf(`
locals {
  category_ids = try(split(",", "%[1]s"), [])
}

resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  count = 2

  workspace_id = "%[2]s"
  name         = format("%[3]s_%%d", count.index) 
}

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  count = 2

  workspace_id     = "%[2]s"
  rule_type        = "CUSTOM"
  name             = format("%[3]s_%%d", count.index)
  secrecy_level_id = huaweicloud_dataarts_security_data_secrecy_level.test[count.index].id
  category_id      = local.category_ids[count.index]
  method           = "NONE"
}
`, acceptance.HW_DATAARTS_SECURITY_DATA_CATEGORY_IDS, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityDataRecognitionRuleGroup_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule_group" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"
  description  = "Created by terraform"
  rule_ids     = slice(huaweicloud_dataarts_security_data_recognition_rule.test[*].id, 0, 1)
}
`, testAccSecurityDataRecognitionRuleGroup_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityDataRecognitionRuleGroup_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule_group" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s_update"
  rule_ids     = slice(huaweicloud_dataarts_security_data_recognition_rule.test[*].id, 1, 2)
}
`, testAccSecurityDataRecognitionRuleGroup_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
