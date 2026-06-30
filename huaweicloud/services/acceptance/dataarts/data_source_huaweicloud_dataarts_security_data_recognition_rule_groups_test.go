package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityDataRecognitionRuleGroups_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_security_data_recognition_rule_groups.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_security_data_recognition_rule_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byCreator   = "data.huaweicloud_dataarts_security_data_recognition_rule_groups.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsSecurityDataCategoryIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSecurityDataRecognitionRuleGroups_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Security data recognition rule groups"),
			},
			{
				Config: testAccDataSecurityDataRecognitionRuleGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "groups.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "groups.0.id"),
					resource.TestCheckResourceAttrSet(all, "groups.0.project_id"),
					resource.TestCheckResourceAttrSet(all, "groups.0.created_by"),
					resource.TestCheckResourceAttrSet(all, "groups.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "groups.0.rules.0.id"),
					resource.TestCheckResourceAttr(all, "groups.0.rules.0.type", "CUSTOM"),

					// Filter by name
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttr(byName, "groups.#", "1"),
					resource.TestCheckResourceAttr(byName, "groups.0.rules.#", "2"),

					// Filter by creator
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityDataRecognitionRuleGroups_nonExistentWorkspace() string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_security_data_recognition_rule_groups" "test" {
  workspace_id = "%[1]s"
}
`, randomUUID.String())
}

func testAccDataSecurityDataRecognitionRuleGroups_basic_base(name string) string {
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

resource "huaweicloud_dataarts_security_data_recognition_rule_group" "test" {
  workspace_id = "%[2]s"
  name         = "%[3]s"
  description  = "Created by terraform"
  rule_ids     = huaweicloud_dataarts_security_data_recognition_rule.test[*].id
}
`, acceptance.HW_DATAARTS_SECURITY_DATA_CATEGORY_IDS, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccDataSecurityDataRecognitionRuleGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all groups and without any filter
data "huaweicloud_dataarts_security_data_recognition_rule_groups" "all" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule_group.test,
  ]

  workspace_id = "%[2]s"
}

# Filter by name
locals {
  group_name = huaweicloud_dataarts_security_data_recognition_rule_group.test.name
}

data "huaweicloud_dataarts_security_data_recognition_rule_groups" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule_group.test,
  ]

  workspace_id = "%[2]s"
  name         = local.group_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_security_data_recognition_rule_groups.filter_by_name.groups[*].name :
      v == local.group_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by creator
locals {
  creator = huaweicloud_dataarts_security_data_recognition_rule_group.test.created_by
}

data "huaweicloud_dataarts_security_data_recognition_rule_groups" "filter_by_creator" {
  depends_on = [
    huaweicloud_dataarts_security_data_recognition_rule_group.test,
  ]

  workspace_id = "%[2]s"
  creator      = local.creator
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_security_data_recognition_rule_groups.filter_by_creator.groups[*].created_by :
      v == local.creator
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}
`, testAccDataSecurityDataRecognitionRuleGroups_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
