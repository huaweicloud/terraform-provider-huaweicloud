package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityDataCategories_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_data_categories.test"
		dc  = acceptance.InitDataSourceCheck(all)
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
				Config: testAccDataSecurityDataCategories_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "categories.#", regexp.MustCompile(`^[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "categories.0.category_id"),
					resource.TestCheckResourceAttrSet(all, "categories.0.category_name"),
					resource.TestCheckResourceAttrSet(all, "categories.0.category_level"),
					resource.TestCheckResourceAttrSet(all, "categories.0.root_id"),
					resource.TestCheckResourceAttrSet(all, "categories.0.parent_id"),
					resource.TestCheckResourceAttrSet(all, "categories.0.category_path"),
					resource.TestCheckResourceAttrSet(all, "categories.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "categories.0.create_by"),
					resource.TestMatchResourceAttr(all, "categories.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_rules_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityDataCategories_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_data_secrecy_level" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
}

locals {
  category_id = try(split(",", "%[3]s")[0], null)
}

resource "huaweicloud_dataarts_security_data_recognition_rule" "test" {
  workspace_id       = "%[1]s"
  rule_type          = "CUSTOM"
  name               = "%[2]s"
  secrecy_level_id   = huaweicloud_dataarts_security_data_secrecy_level.test.id
  category_id        = local.category_id
  method             = "REGULAR"
  content_expression = "^male$|^female&"
  column_expression  = "phoneNumber|email"
  comment_expression = ".*comment*."
  description        = "Created by terraform script"
  enable             = true
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_SECURITY_DATA_CATEGORY_IDS)
}

func testAccDataSecurityDataCategories_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_dataarts_security_data_categories" "test" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_security_data_recognition_rule.test]
}

locals {
  associated_rule = try(element([for v in data.huaweicloud_dataarts_security_data_categories.test.categories : v.rules
  if v.category_id == local.category_id][0], 0), null)
}

output "is_rules_set_and_valid" {
  value = try(local.associated_rule.id == huaweicloud_dataarts_security_data_recognition_rule.test.id &&
    local.associated_rule.name == huaweicloud_dataarts_security_data_recognition_rule.test.name &&
    local.associated_rule.enable == huaweicloud_dataarts_security_data_recognition_rule.test.enable &&
    local.associated_rule.method == huaweicloud_dataarts_security_data_recognition_rule.test.method &&
    local.associated_rule.content_expression == huaweicloud_dataarts_security_data_recognition_rule.test.content_expression &&
    local.associated_rule.column_expression == huaweicloud_dataarts_security_data_recognition_rule.test.column_expression &&
    local.associated_rule.comment_expression == huaweicloud_dataarts_security_data_recognition_rule.test.comment_expression &&
  local.associated_rule.description == huaweicloud_dataarts_security_data_recognition_rule.test.description, false)
}
`, testAccDataSecurityDataCategories_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
