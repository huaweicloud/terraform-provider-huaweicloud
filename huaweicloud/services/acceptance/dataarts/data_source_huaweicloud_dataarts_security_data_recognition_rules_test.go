package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityDataRecognitionRules_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_data_recognition_rules.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		bySecrecyLevel   = "data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_secrecy_level"
		dcBySecrecyLevel = acceptance.InitDataSourceCheck(bySecrecyLevel)

		byCreator   = "data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)

		byEnable   = "data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_enable"
		dcByEnable = acceptance.InitDataSourceCheck(byEnable)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsSecurityDataCategoryIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityDataRecognitionRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttr(byName, "rules.#", "1"),
					resource.TestCheckResourceAttr(byName, "rules.0.name", rName),
					resource.TestCheckResourceAttr(byName, "rules.0.type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(byName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(byName, "rules.0.enable"),
					resource.TestCheckResourceAttrSet(byName, "rules.0.created_at"),
					resource.TestCheckResourceAttrSet(byName, "rules.0.created_by"),
					resource.TestCheckResourceAttrSet(byName, "rules.0.updated_at"),
					resource.TestCheckResourceAttrSet(byName, "rules.0.updated_by"),

					// Filter by secrecy_level
					dcBySecrecyLevel.CheckResourceExists(),
					resource.TestCheckOutput("is_secrecy_level_filter_useful", "true"),
					resource.TestMatchResourceAttr(bySecrecyLevel, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(bySecrecyLevel, "rules.0.secrecy_level"),

					// Filter by creator
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
					resource.TestMatchResourceAttr(byCreator, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(byCreator, "rules.0.created_by"),

					// Filter by enable
					dcByEnable.CheckResourceExists(),
					resource.TestCheckOutput("is_enable_filter_useful", "true"),
					resource.TestMatchResourceAttr(byEnable, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(byEnable, "rules.0.enable", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityDataRecognitionRules_basic_base(name string) string {
	return fmt.Sprintf(`
locals {
  category_ids = try(split(",", "%[1]s"), [])
}

resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"
}

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id     = "%[2]s"
  rule_type        = "CUSTOM"
  name             = "%[3]s"
  secrecy_level_id = huaweicloud_dataarts_security_data_secrecy_level.test.id
  category_id      = try(local.category_ids[0], null)
  method           = "NONE"
}
`, acceptance.HW_DATAARTS_SECURITY_DATA_CATEGORY_IDS, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataSecurityDataRecognitionRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all rules and without any filter
data "huaweicloud_dataarts_security_data_recognition_rules" "all" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule.test,
  ]

  workspace_id = "%[2]s"
}

# Filter by name
locals {
  rule_name = huaweicloud_dataarts_security_data_recognition_rule.test.name
}

data "huaweicloud_dataarts_security_data_recognition_rules" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule.test,
  ]

  workspace_id = "%[2]s"
  name         = local.rule_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_name.rules[*].name :
      v == local.rule_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by secrecy_level
locals {
  secrecy_level = huaweicloud_dataarts_security_data_recognition_rule.test.secrecy_level
}

data "huaweicloud_dataarts_security_data_recognition_rules" "filter_by_secrecy_level" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule.test,
  ]

  workspace_id  = "%[2]s"
  secrecy_level = huaweicloud_dataarts_security_data_recognition_rule.test.secrecy_level
}

locals {
  secrecy_level_filter_result = [
    for v in data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_secrecy_level.rules[*].secrecy_level :
      v == local.secrecy_level
  ]
}

output "is_secrecy_level_filter_useful" {
  value = length(local.secrecy_level_filter_result) > 0 && alltrue(local.secrecy_level_filter_result)
}

# Filter by creator
locals {
  creator = huaweicloud_dataarts_security_data_recognition_rule.test.created_by
}

data "huaweicloud_dataarts_security_data_recognition_rules" "filter_by_creator" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule.test,
  ]

  workspace_id = "%[2]s"
  creator      = huaweicloud_dataarts_security_data_recognition_rule.test.created_by
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_creator.rules[*].created_by :
    v == local.creator
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}

# Filter by enable
locals {
  is_rule_enable = true
}

data "huaweicloud_dataarts_security_data_recognition_rules" "filter_by_enable" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule.test,
  ]

  workspace_id = "%[2]s"
  enable       = local.is_rule_enable
}

locals {
  enable_filter_result = [
    for v in data.huaweicloud_dataarts_security_data_recognition_rules.filter_by_enable.rules[*].enable :
      v == local.is_rule_enable
  ]
}

output "is_enable_filter_useful" {
  value = length(local.enable_filter_result) > 0 && alltrue(local.enable_filter_result)
}
`, testAccDataSecurityDataRecognitionRules_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
