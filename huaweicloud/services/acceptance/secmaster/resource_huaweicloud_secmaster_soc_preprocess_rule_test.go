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

func getResourceSocPreprocessRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetSocPreprocessRule(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceSocPreprocessRule_basic(t *testing.T) {
	var (
		rName = "huaweicloud_secmaster_soc_preprocess_rule.test"

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceSocPreprocessRuleFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterMappingId(t)
			acceptance.TestAccPreCheckSecMasterMapperId(t)
			acceptance.TestAccPreCheckSecMasterMapperTypeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSocPreprocessRule_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "mapping_id", acceptance.HW_SECMASTER_MAPPING_ID),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.mapper_id", acceptance.HW_SECMASTER_MAPPER_ID),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.mapper_type_id", acceptance.HW_SECMASTER_MAPPER_TYPE_ID),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.action", "drop"),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.expression", "expression_content"),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.name", "test_name"),
					resource.TestCheckResourceAttrSet(rName, "data.#"),
				),
			},
			{
				Config: testAccSocPreprocessRule_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.mapper_id", acceptance.HW_SECMASTER_MAPPER_ID),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.mapper_type_id", acceptance.HW_SECMASTER_MAPPER_TYPE_ID),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.name", "update_name"),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.action", "drop"),
					resource.TestCheckResourceAttr(rName, "preprocess_rules.0.expression", "update_expression"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSocPreprocessRuleImportStateFunc(rName),
			},
		},
	})
}

func testAccSocPreprocessRule_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_soc_preprocess_rule" "test" {
  workspace_id = "%[1]s"
  mapping_id   = "%[2]s"

  preprocess_rules {
    name           = "test_name"
    mapper_id      = "%[3]s"
    mapper_type_id = "%[4]s"
    action         = "drop"
    expression     = "expression_content"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_MAPPING_ID,
		acceptance.HW_SECMASTER_MAPPER_ID, acceptance.HW_SECMASTER_MAPPER_TYPE_ID)
}

func testAccSocPreprocessRule_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_soc_preprocess_rule" "test" {
  workspace_id = "%[1]s"
  mapping_id   = "%[2]s"

  preprocess_rules {
    name           = "update_name"
    mapper_id      = "%[3]s"
    mapper_type_id = "%[4]s"
    action         = "drop"
    expression     = "update_expression"
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_MAPPING_ID,
		acceptance.HW_SECMASTER_MAPPER_ID, acceptance.HW_SECMASTER_MAPPER_TYPE_ID)
}

func testAccSocPreprocessRuleImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, mappingId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		mappingId = rs.Primary.ID

		if workspaceId == "" || mappingId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<mapping_id>', but got '%s/%s'",
				workspaceId, mappingId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, mappingId), nil
	}
}
