package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrchestrationRuleAssociatedApis_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_apig_orchestration_rule_associated_apis.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byApiId   = "data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_api_id"
		dcByApiId = acceptance.InitDataSourceCheck(byApiId)

		byNotFoundApiId   = "data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_not_found_api_id"
		dcByNotFoundApiId = acceptance.InitDataSourceCheck(byNotFoundApiId)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceOrchestrationRuleAssociatedApis_basic_step1(),
				ExpectError: regexp.MustCompile(`The instance does not exist`),
			},
			{
				Config:      testAccDataSourceOrchestrationRuleAssociatedApis_basic_step2(),
				ExpectError: regexp.MustCompile(`The orchestrations does not exist`),
			},
			{
				Config: testAccDataSourceOrchestrationRuleAssociatedApis_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "apis.#", regexp.MustCompile(`[1-9]\d*`)),
					// Filter by API ID.
					dcByApiId.CheckResourceExists(),
					resource.TestCheckOutput("is_api_id_filter_useful", "true"),
					dcByNotFoundApiId.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_api_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceOrchestrationRuleAssociatedApis_basic_step1() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_apig_orchestration_rule_associated_apis" "incorrect_instance_id" {
  instance_id = "%[1]s"
  rule_id     = "%[1]s"
}
`, randUUID)
}

func testAccDataSourceOrchestrationRuleAssociatedApis_basic_step2() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_apig_orchestration_rule_associated_apis" "incorrect_orchestration_rule_id" {
  instance_id = "%[1]s"
  rule_id     = "%[2]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, randUUID)
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
    "mapped_param_name": "standard-param",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}

resource "huaweicloud_apig_api" "test" {
  count = 2

  instance_id             = "%[1]s"
  group_id                = huaweicloud_apig_group.test.id
  name                    = format("%[2]s_%%d", count.index)
  type                    = "Private"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = format("/orchestration/test/%%d", count.index)
  security_authentication = "NONE"

  request_params {
    name     = "X-Service-Name"
    type     = "STRING"
    location = "HEADER"
    maximum  = 30
    minimum  = 5

    orchestrations = [
      huaweicloud_apig_orchestration_rule.test.id,
    ]
  }

  backend_params {
    type              = "REQUEST"
    name              = "ServiceName"
    location          = "HEADER"
    value             = "X-Service-Name"
    system_param_type = "backend"
  }

  mock {
    status_code   = 201
    response      = "{'message':'hello world'}"
  }
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceOrchestrationRuleAssociatedApis_basic_step3(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_orchestration_rule_associated_apis" "all" {
  depends_on = [
    huaweicloud_apig_api.test,
  ]

  instance_id = "%[2]s"
  rule_id     = huaweicloud_apig_orchestration_rule.test.id
}

# Filter by ID
locals {
  api_id = huaweicloud_apig_api.test[0].id
}

data "huaweicloud_apig_orchestration_rule_associated_apis" "filter_by_api_id" {
  # Need to be executed after all resources are created (excluding implicit dependencies).
  depends_on = [
    huaweicloud_apig_api.test,
  ]

  instance_id = "%[2]s"
  rule_id     = huaweicloud_apig_orchestration_rule.test.id
  api_id      = local.api_id
}

locals {
  api_id_filter_result = [
    for v in data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_api_id.apis[*].api_id : v == local.api_id
  ]
}

output "is_api_id_filter_useful" {
  value = length(local.api_id_filter_result) > 0 && alltrue(local.api_id_filter_result)
}

# Filter by not found ID
data "huaweicloud_apig_orchestration_rule_associated_apis" "filter_by_not_found_api_id" {
  # Need to be executed after all resources are created (excluding implicit dependencies).
  depends_on = [
    huaweicloud_apig_api.test,
  ]

  instance_id = "%[2]s"
  rule_id     = huaweicloud_apig_orchestration_rule.test.id
  api_id      = "%[3]s"
}

output "is_not_found_api_id_filter_useful" {
  value = length(data.huaweicloud_apig_orchestration_rule_associated_apis.filter_by_not_found_api_id.apis) == 0
}
`, testAccDataSourceOrchestrationRuleAssociatedApis_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, randUUID)
}
