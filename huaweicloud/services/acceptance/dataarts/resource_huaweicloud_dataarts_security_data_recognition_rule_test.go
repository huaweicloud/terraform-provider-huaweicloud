package dataarts

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

func getSecurityRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getRuleClient, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio V1 client: %s", err)
	}

	getRuleHttpUrl := "v1/{project_id}/security/data-classification/rule/{id}"

	getRulePath := getRuleClient.Endpoint + getRuleHttpUrl
	getRulePath = strings.ReplaceAll(getRulePath, "{project_id}", getRuleClient.ProjectID)
	getRulePath = strings.ReplaceAll(getRulePath, "{id}", state.Primary.ID)

	getRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}
	getRuleResp, err := getRuleClient.Request("GET", getRulePath, &getRuleOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getRuleResp)
}

func TestAccResourceSecurityRule_custom(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_security_data_recognition_rule.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSecurityRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsDataClassificationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityRule_custom(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "method", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "category_id", acceptance.HW_DATAARTS_CATEGORY_ID),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttrSet(resourceName, "enable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccSecurityRule_custom_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "method", "REGULAR"),
					resource.TestCheckResourceAttr(resourceName, "category_id", acceptance.HW_DATAARTS_CATEGORY_ID_UPDATE),
					resource.TestCheckResourceAttr(resourceName, "content_expression", "^male$|^female&"),
					resource.TestCheckResourceAttr(resourceName, "column_expression", "phoneNumber|email"),
					resource.TestCheckResourceAttr(resourceName, "comment_expression", ".*comment*."),
					resource.TestCheckResourceAttr(resourceName, "description", "rule_description_custom_update1"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttrSet(resourceName, "enable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccSecurityRule_custom_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "method", "REGULAR"),
					resource.TestCheckResourceAttr(resourceName, "category_id", acceptance.HW_DATAARTS_CATEGORY_ID),
					resource.TestCheckResourceAttr(resourceName, "content_expression", "^boy$|^girl&"),
					resource.TestCheckResourceAttr(resourceName, "column_expression", "sex|gender"),
					resource.TestCheckResourceAttr(resourceName, "comment_expression", ".*commentUpdate*."),
					resource.TestCheckResourceAttr(resourceName, "description", "rule_description_custom_update2"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttrSet(resourceName, "enable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testDataArtsStudioImportState(resourceName),
				ImportStateVerifyIgnore: []string{"secrecy_level_id"},
			},
		},
	})
}

func TestAccResourceSecurityRule_builtin(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_security_data_recognition_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSecurityRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsDataClassificationID(t)
			acceptance.TestAccPreCheckDataArtsBuiltinRule(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityRule_builtin(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_DATAARTS_BUILTIN_RULE_NAME),
					resource.TestCheckResourceAttr(resourceName, "builtin_rule_id", acceptance.HW_DATAARTS_BUILTIN_RULE_ID),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "BUILTIN"),
					resource.TestCheckResourceAttr(resourceName, "category_id", acceptance.HW_DATAARTS_CATEGORY_ID),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttrSet(resourceName, "enable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccSecurityRule_builtin_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "BUILTIN"),
					resource.TestCheckResourceAttr(resourceName, "builtin_rule_id", acceptance.HW_DATAARTS_BUILTIN_RULE_ID),
					resource.TestCheckResourceAttr(resourceName, "secrecy_level_id", acceptance.HW_DATAARTS_SECRECY_LEVEL_ID_UPDATE),
					resource.TestCheckResourceAttr(resourceName, "category_id", acceptance.HW_DATAARTS_CATEGORY_ID_UPDATE),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttrSet(resourceName, "enable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
		},
	})
}

func testAccSecurityRule_custom(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[1]s"
  rule_type        = "CUSTOM"
  name             = "%[2]s"
  secrecy_level_id = "%[3]s"
  category_id      = "%[4]s"
  method           = "NONE"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_SECRECY_LEVEL_ID, acceptance.HW_DATAARTS_CATEGORY_ID)
}

func testAccSecurityRule_custom_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id       = "%[1]s"
  rule_type          = "CUSTOM"
  name               = "%[2]s"
  secrecy_level_id   = "%[3]s"
  category_id        = "%[4]s"
  method             = "REGULAR"
  content_expression = "^male$|^female&"
  column_expression  = "phoneNumber|email"
  comment_expression = ".*comment*."
  description        = "rule_description_custom_update1"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_SECRECY_LEVEL_ID_UPDATE, acceptance.HW_DATAARTS_CATEGORY_ID_UPDATE)
}

func testAccSecurityRule_custom_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id       = "%[1]s"
  rule_type          = "CUSTOM"
  name               = "%[2]s"
  secrecy_level_id   = "%[3]s"
  category_id        = "%[4]s"
  method             = "REGULAR"
  content_expression = "^boy$|^girl&"
  column_expression  = "sex|gender"
  comment_expression = ".*commentUpdate*."
  description        = "rule_description_custom_update2"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_SECRECY_LEVEL_ID, acceptance.HW_DATAARTS_CATEGORY_ID)
}

func testAccSecurityRule_builtin() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[1]s"
  rule_type        = "BUILTIN"
  name             = "%[2]s"
  builtin_rule_id  = "%[3]s"
  secrecy_level_id = "%[4]s"
  category_id      = "%[5]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_BUILTIN_RULE_NAME, acceptance.HW_DATAARTS_BUILTIN_RULE_ID,
		acceptance.HW_DATAARTS_SECRECY_LEVEL_ID, acceptance.HW_DATAARTS_CATEGORY_ID)
}

func testAccSecurityRule_builtin_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[1]s"
  rule_type        = "BUILTIN"
  name             = "%[2]s"
  builtin_rule_id  = "%[3]s"
  secrecy_level_id = "%[4]s"
  category_id      = "%[5]s"
  description      = "rule_description_builtin_update"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_BUILTIN_RULE_NAME, acceptance.HW_DATAARTS_BUILTIN_RULE_ID,
		acceptance.HW_DATAARTS_SECRECY_LEVEL_ID_UPDATE, acceptance.HW_DATAARTS_CATEGORY_ID_UPDATE)
}
