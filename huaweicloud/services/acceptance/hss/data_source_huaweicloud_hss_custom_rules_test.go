package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCustomRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_custom_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rule_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rule_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.auto_block"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.hash_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.is_all_host"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.update_time"),

					resource.TestCheckOutput("is_rule_id_filter_useful", "true"),
					resource.TestCheckOutput("is_rule_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCustomRules_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_custom_rule" "test" {
  rule_name   = "%[1]s"
  is_all_host = true
  agent_ids   = ["af08124fd77581dcd2a5f7cdaa208c285af0a1480f7ca9b5258193c021ca2637"] // mock data
  rule_status = 0

  custom_rule_value_info {
    auto_block  = 1
    hash_type   = "sha1"
    rule_type   = "black_hash"
    rule_values = ["08a7baa28dd268f8a12bc1f6fd95869321fe51144c5bf3321a6f6305edcd5245"] // mock data
  }
}
`, name)
}

func testAccDataSourceCustomRules_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_custom_rules" "test" {
  depends_on = [huaweicloud_hss_custom_rule.test]
}

# Filter using rule_id.
locals {
  rule_id = data.huaweicloud_hss_custom_rules.test.data_list[0].rule_id
}

data "huaweicloud_hss_custom_rules" "rule_id_filter" {
  rule_id = local.rule_id
}

output "is_rule_id_filter_useful" {
  value = length(data.huaweicloud_hss_custom_rules.rule_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_custom_rules.rule_id_filter.data_list[*].rule_id : v == local.rule_id]
  )
}

# Filter using rule_name.
locals {
  rule_name = data.huaweicloud_hss_custom_rules.test.data_list[0].rule_name
}

data "huaweicloud_hss_custom_rules" "rule_name_filter" {
  rule_name = local.rule_name
}

output "is_rule_name_filter_useful" {
  value = length(data.huaweicloud_hss_custom_rules.rule_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_custom_rules.rule_name_filter.data_list[*].rule_name : v == local.rule_name]
  )
}

# Filter using non existent rule_name.
data "huaweicloud_hss_custom_rules" "not_found" {
  rule_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_custom_rules.not_found.data_list) == 0
}
`, testAccDataSourceCustomRules_base(name))
}
