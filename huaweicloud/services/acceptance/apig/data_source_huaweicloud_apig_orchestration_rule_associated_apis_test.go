package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrchestrationRuleAssociatedApis_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_apig_orchestration_rule_associated_apis.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byApiId   = "data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_api_id"
		dcByApiId = acceptance.InitDataSourceCheck(byApiId)

		byApiName   = "data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_api_name"
		dcByApiName = acceptance.InitDataSourceCheck(byApiName)

		apiNameNotFound   = "data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_not_found_api_name"
		dcApiNameNotFound = acceptance.InitDataSourceCheck(apiNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrchestrationRuleAssociatedApis_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "apis.#", regexp.MustCompile(`[1-9]\d*`)),
					dcByApiId.CheckResourceExists(),
					resource.TestCheckOutput("api_id_filter_is_useful", "true"),
					dcByApiName.CheckResourceExists(),
					resource.TestCheckOutput("api_name_filter_is_useful", "true"),
					dcApiNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(byApiId, "apis.0.api_id"),
					resource.TestCheckResourceAttrSet(byApiId, "apis.0.api_name"),
					resource.TestCheckResourceAttrSet(byApiId, "apis.0.req_uri"),
					resource.TestCheckResourceAttrSet(byApiId, "apis.0.req_method"),
					resource.TestCheckResourceAttrSet(byApiId, "apis.0.auth_type"),
					resource.TestCheckResourceAttrSet(byApiId, "apis.0.match_mode"),
					resource.TestCheckResourceAttrPair(byApiId, "apis.0.group_id", "huaweicloud_apig_group.test", "id"),
					resource.TestCheckResourceAttrPair(byApiId, "apis.0.group_name", "huaweicloud_apig_group.test", "name"),
					resource.TestMatchResourceAttr(byApiId, "apis.0.attached_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceOrchestrationRuleAssociatedApis_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  strategy    = "hash"

  mapped_param = jsonencode({
    "mapped_param_name" : "standard-param",
    "mapped_param_type" : "string",
    "mapped_param_location" : "header"
  })
}

resource "huaweicloud_apig_api" "test" {
  count = 2

  instance_id      = "%[1]s"
  group_id         = huaweicloud_apig_group.test.id
  type             = "Private"
  name             = "%[2]s${count.index}"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = "/test${count.index}/users"

  request_params {
    name           = "X-Service-Name"
    type           = "STRING"
    location       = "HEADER"
    orchestrations = [huaweicloud_apig_orchestration_rule.test.id]
  }

  mock {
    status_code = 200
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceOrchestrationRuleAssociatedApis_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_orchestration_rule_associated_apis" "test" {
  instance_id      = "%[2]s"
  orchestration_id = huaweicloud_apig_orchestration_rule.test.id

  depends_on = [huaweicloud_apig_api.test]
}

locals {
  api_id   = huaweicloud_apig_api.test[0].id
  api_name = huaweicloud_apig_api.test[0].name
}

data "huaweicloud_apig_orchestration_rule_associated_apis" "filter_by_api_id" {
  instance_id      = "%[2]s"
  orchestration_id = huaweicloud_apig_orchestration_rule.test.id
  api_id           = local.api_id
}

locals {
  api_id_filter_result = [for v in data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_api_id.apis[*].api_id : v == local.api_id]
}

output "api_id_filter_is_useful" {
  value = length(local.api_id_filter_result) > 0 && alltrue(local.api_id_filter_result)
}

# Filter by API name.
data "huaweicloud_apig_orchestration_rule_associated_apis" "filter_by_api_name" {
  instance_id      = "%[2]s"
  orchestration_id = huaweicloud_apig_orchestration_rule.test.id
  api_name         = local.api_name

  depends_on = [huaweicloud_apig_api.test]
}

locals {
  api_name_filter_result = [for v in data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_api_name.apis[*].api_name :
  strcontains(v, local.api_name)]
}

output "api_name_filter_is_useful" {
  value = length(local.api_name_filter_result) > 0 && alltrue(local.api_name_filter_result)
}

# Filter by non-existent API name.
data "huaweicloud_apig_orchestration_rule_associated_apis" "filter_by_not_found_api_name" {
  instance_id      = "%[2]s"
  orchestration_id = huaweicloud_apig_orchestration_rule.test.id
  api_name         = "not_found_name"

  depends_on = [huaweicloud_apig_api.test]
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_not_found_api_name.apis) == 0
}
`, testAccDataSourceOrchestrationRuleAssociatedApis_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
