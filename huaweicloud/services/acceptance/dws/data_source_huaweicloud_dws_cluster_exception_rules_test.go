package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterExceptionRules_basic(t *testing.T) {
	var (
		all   = "data.huaweicloud_dws_cluster_exception_rules.all"
		dcAll = acceptance.InitDataSourceCheck(all)

		byNameForExactMatch   = "data.huaweicloud_dws_cluster_exception_rules.filter_by_rule_name_for_exact_match"
		dcByNameForExactMatch = acceptance.InitDataSourceCheck(byNameForExactMatch)

		byNameForFuzzyMatch   = "data.huaweicloud_dws_cluster_exception_rules.filter_by_rule_name_for_fuzzy_match"
		dcByNameForFuzzyMatch = acceptance.InitDataSourceCheck(byNameForFuzzyMatch)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataClusterExceptionRules_nonExistentCluster(),
				ExpectError: regexp.MustCompile(`error querying cluster exception rules`),
			},
			{
				Config: testAccDataClusterExceptionRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`^[0-9]+$`)),
					// Filter by rule name.
					dcByNameForExactMatch.CheckResourceExists(),
					resource.TestCheckOutput("is_rule_name_for_exact_match_filter_useful", "true"),
					dcByNameForFuzzyMatch.CheckResourceExists(),
					resource.TestCheckOutput("is_rule_name_for_fuzzy_match_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataClusterExceptionRules_nonExistentCluster() string {
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_exception_rules" "nonexistent_cluster" {
  cluster_id = "%[1]s"
}
`, randomUUID)
}

func testAccDataClusterExceptionRules_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "exception_rule_configurations" {
  description = "The configurations of the exception rule."
  type        = list(object({
    key   = string
    value = string
  }))

  default = [
    {
      key   = "action"
      value = "penalty"
    },
    {
      key   = "blocktime"
      value = "300"
    },
    {
      key   = "elapsedtime"
      value = "400"
    },
    {
        key   = "allcputime"
        value = "500"
    },
  ]
}

resource "huaweicloud_dws_cluster_exception_rule" "test" {
  count = 2

  cluster_id = "%[1]s"
  name       = format("%[2]s_%%d", count.index)

  dynamic "configurations" {
    for_each = var.exception_rule_configurations

    content {
      key   = configurations.value.key
      value = configurations.value.value
    }
  }
}

# Query all exception rules in the cluster
data "huaweicloud_dws_cluster_exception_rules" "all" {
  depends_on = [
    huaweicloud_dws_cluster_exception_rule.test,
  ]

  cluster_id = "%[1]s"
}

# Filter the exception rules with a specified name in the cluster (fuzzy matching)
locals {
  rule_name_for_exact_match = huaweicloud_dws_cluster_exception_rule.test[0].name
  rule_name_for_fuzzy_match = "%[2]s"
}

data "huaweicloud_dws_cluster_exception_rules" "filter_by_rule_name_for_exact_match" {
  # The behavior of parameter 'rule_name' of the exception rule resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_dws_cluster_exception_rule.test,
  ]

  cluster_id = "%[1]s"
  rule_name  = local.rule_name_for_exact_match
}

locals {
  rule_name_for_exact_match_filter_result = [
    for v in data.huaweicloud_dws_cluster_exception_rules.filter_by_rule_name_for_exact_match.rules[*].name : v == local.rule_name_for_exact_match
  ]
}

output "is_rule_name_for_exact_match_filter_useful" {
  value = length(local.rule_name_for_exact_match_filter_result) > 0 && alltrue(local.rule_name_for_exact_match_filter_result)
}

data "huaweicloud_dws_cluster_exception_rules" "filter_by_rule_name_for_fuzzy_match" {
  # The input value of parameter 'rule_name' is a fixed value, so we need to depend on the exception rule resource.
  depends_on = [
    huaweicloud_dws_cluster_exception_rule.test,
  ]

  cluster_id = "%[1]s"
  rule_name  = local.rule_name_for_fuzzy_match
}

locals {
  rule_name_for_fuzzy_match_filter_result = [
    for v in data.huaweicloud_dws_cluster_exception_rules.filter_by_rule_name_for_fuzzy_match.rules[*].name :
      length(regexall(format(".*%%s.*", local.rule_name_for_fuzzy_match), v)) > 0
  ]
}

output "is_rule_name_for_fuzzy_match_filter_useful" {
  value = length(local.rule_name_for_fuzzy_match_filter_result) > 0 && alltrue(local.rule_name_for_fuzzy_match_filter_result)
}
`, acceptance.HW_DWS_CLUSTER_ID, name)
}
