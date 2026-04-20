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

func getSecurityDataRecognitionRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccSecurityDataRecognitionRule_custom(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_security_data_recognition_rule.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getSecurityDataRecognitionRuleResourceFunc)

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
				Config: testAccSecurityDataRecognitionRule_custom_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "method", "NONE"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccSecurityDataRecognitionRule_custom_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "method", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttr(resourceName, "content_expression", "^male$|^female&"),
					resource.TestCheckResourceAttr(resourceName, "column_expression", "phoneNumber|email"),
					resource.TestCheckResourceAttr(resourceName, "comment_expression", ".*comment*."),
					resource.TestCheckResourceAttr(resourceName, "description", "rule_description_custom_update1"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttr(resourceName, "enable", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccSecurityDataRecognitionRule_custom_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "CUSTOM"),
					resource.TestCheckResourceAttr(resourceName, "method", "REGULAR"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
					resource.TestCheckResourceAttr(resourceName, "content_expression", "^boy$|^girl&"),
					resource.TestCheckResourceAttr(resourceName, "column_expression", "sex|gender"),
					resource.TestCheckResourceAttr(resourceName, "comment_expression", ".*commentUpdate*."),
					resource.TestCheckResourceAttr(resourceName, "description", "rule_description_custom_update2"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level"),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_num"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
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

func testAccSecurityDataRecognitionRule_base(name string) string {
	return fmt.Sprintf(`
locals {
  category_ids = try(split(",", "%[1]s"), [])
}

resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  count = 2

  workspace_id = "%[2]s"
  name         = format("%[3]s_%%d", count.index) 
}
`, acceptance.HW_DATAARTS_SECURITY_DATA_CATEGORY_IDS, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityDataRecognitionRule_custom_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[2]s"
  rule_type        = "CUSTOM"
  name             = "%[3]s"
  secrecy_level_id = huaweicloud_dataarts_security_data_secrecy_level.test[0].id
  category_id      = local.category_ids[0]
  method           = "NONE"
}
`, testAccSecurityDataRecognitionRule_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityDataRecognitionRule_custom_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id       = "%[2]s"
  rule_type          = "CUSTOM"
  name               = "%[3]s"
  secrecy_level_id   = huaweicloud_dataarts_security_data_secrecy_level.test[1].id
  category_id        = local.category_ids[1]
  method             = "REGULAR"
  content_expression = "^male$|^female&"
  column_expression  = "phoneNumber|email"
  comment_expression = ".*comment*."
  description        = "rule_description_custom_update1"
  enable             = false
}
`, testAccSecurityDataRecognitionRule_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccSecurityDataRecognitionRule_custom_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id       = "%[2]s"
  rule_type          = "CUSTOM"
  name               = "%[3]s"
  secrecy_level_id   = huaweicloud_dataarts_security_data_secrecy_level.test[0].id
  category_id        = local.category_ids[0]
  method             = "REGULAR"
  content_expression = "^boy$|^girl&"
  column_expression  = "sex|gender"
  comment_expression = ".*commentUpdate*."
  description        = "rule_description_custom_update2"
  enable             = true
}
`, testAccSecurityDataRecognitionRule_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func TestAccSecurityDataRecognitionRule_builtin(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_dataarts_security_data_recognition_rule.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getSecurityDataRecognitionRuleResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsSecurityDataCategoryIds(t, 2)
			acceptance.TestAccPreCheckDataArtsBuiltinRule(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityDataRecognitionRule_builtin_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "BUILTIN"),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_NAME),
					resource.TestCheckResourceAttr(resourceName, "builtin_rule_id", acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_ID),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
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
				Config: testAccSecurityDataRecognitionRule_builtin_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "rule_type", "BUILTIN"),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_NAME),
					resource.TestCheckResourceAttr(resourceName, "builtin_rule_id", acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_ID),
					resource.TestCheckResourceAttrSet(resourceName, "secrecy_level_id"),
					resource.TestCheckResourceAttrSet(resourceName, "category_id"),
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

func testAccSecurityDataRecognitionRule_builtin_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[2]s"
  rule_type        = "BUILTIN"
  name             = "%[3]s"
  builtin_rule_id  = "%[4]s"
  secrecy_level_id = huaweicloud_dataarts_security_data_secrecy_level.test[0].id
  category_id      = local.category_ids[0]
}
`, testAccSecurityDataRecognitionRule_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_NAME,
		acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_ID)
}

func testAccSecurityDataRecognitionRule_builtin_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[2]s"
  rule_type        = "BUILTIN"
  name             = "%[3]s"
  builtin_rule_id  = "%[4]s"
  secrecy_level_id = huaweicloud_dataarts_security_data_secrecy_level.test[1].id
  category_id      = local.category_ids[1]
  description      = "rule_description_builtin_update"
}
`, testAccSecurityDataRecognitionRule_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_NAME,
		acceptance.HW_DATAARTS_SECURITY_BUILTIN_RECOGNITION_RULE_ID)
}
