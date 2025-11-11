package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFunctionTriggers_basic(t *testing.T) {
	var (
		base = "huaweicloud_fgs_function_trigger.test"

		all              = "data.huaweicloud_fgs_function_triggers.all"
		dcForAllTriggers = acceptance.InitDataSourceCheck(all)

		byTriggerId           = "data.huaweicloud_fgs_function_triggers.filter_by_trigger_id"
		dcByTriggerId         = acceptance.InitDataSourceCheck(byTriggerId)
		byNotFoundTriggerId   = "data.huaweicloud_fgs_function_triggers.filter_by_not_found_trigger_id"
		dcByNotFoundTriggerId = acceptance.InitDataSourceCheck(byNotFoundTriggerId)

		byType           = "data.huaweicloud_fgs_function_triggers.filter_by_type"
		dcByType         = acceptance.InitDataSourceCheck(byType)
		byNotFoundType   = "data.huaweicloud_fgs_function_triggers.filter_by_not_found_type"
		dcByNotFoundType = acceptance.InitDataSourceCheck(byNotFoundType)

		byStatus           = "data.huaweicloud_fgs_function_triggers.filter_by_status"
		dcByStatus         = acceptance.InitDataSourceCheck(byStatus)
		byNotFoundStatus   = "data.huaweicloud_fgs_function_triggers.filter_by_not_found_status"
		dcByNotFoundStatus = acceptance.InitDataSourceCheck(byNotFoundStatus)

		byStartTime   = "data.huaweicloud_fgs_function_triggers.filter_by_start_time"
		dcByStartTime = acceptance.InitDataSourceCheck(byStartTime)
		byEndTime     = "data.huaweicloud_fgs_function_triggers.filter_by_end_time"
		dcByEndTime   = acceptance.InitDataSourceCheck(byEndTime)

		byTypeWithoutFunctionUrn   = "data.huaweicloud_fgs_function_triggers.filter_by_type_without_function_urn"
		dcByTypeWithoutFunctionUrn = acceptance.InitDataSourceCheck(byTypeWithoutFunctionUrn)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFunctionTriggers_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllTriggers.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "triggers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by trigger ID.
					dcByTriggerId.CheckResourceExists(),
					resource.TestCheckOutput("is_trigger_id_filter_useful", "true"),
					dcByNotFoundTriggerId.CheckResourceExists(),
					resource.TestCheckOutput("trigger_id_not_found_validation_pass", "true"),
					// Filter by trigger type.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByNotFoundType.CheckResourceExists(),
					resource.TestCheckOutput("type_not_found_validation_pass", "true"),
					// Filter by trigger status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByNotFoundStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_not_found_validation_pass", "true"),
					// Filter by trigger start time.
					dcByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),
					// Filter by trigger end time.
					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrPair(byTriggerId, "triggers.0.id", base, "id"),
					resource.TestCheckResourceAttrPair(byTriggerId, "triggers.0.type", base, "type"),
					resource.TestCheckResourceAttrPair(byTriggerId, "triggers.0.event_data", base, "event_data"),
					resource.TestCheckResourceAttrPair(byTriggerId, "triggers.0.status", base, "status"),
					resource.TestMatchResourceAttr(byTriggerId, "triggers.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byTriggerId, "triggers.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Check the consistency of trigger type between the resource and the data source.
					dcByTypeWithoutFunctionUrn.CheckResourceExists(),
					resource.TestCheckOutput("trigger_type_consistency_check_pass", "true"),
					resource.TestCheckResourceAttrSet(byTypeWithoutFunctionUrn, "triggers.0.function_urn"),
				),
			},
		},
	})
}

func testAccDataFunctionTriggers_basic() string {
	name := acceptance.RandomAccResourceName()
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
variable "function_code_content" {
  type    = string
  default = <<EOT
def main():  
    print("Hello, World!")  

if __name__ == "__main__":  
    main()
EOT
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  handler     = "index.handler"
  app         = "default"
  description = "fuction test"
  code_type   = "inline"
  func_code   = base64encode(var.function_code_content)
}

resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  event_data   = jsonencode({
    name           = "%[1]s_rate",
    schedule_type  = "Rate",
    sync_execution = false,
    user_event     = "Created by terraform script",
    schedule       = "3m"
  })
}

# Without any filter parameter.
data "huaweicloud_fgs_function_triggers" "all" {
  depends_on = [
    huaweicloud_fgs_function_trigger.test
  ]

  function_urn = huaweicloud_fgs_function.test.urn
}

# Filter by trigger ID.
locals {
  trigger_id = huaweicloud_fgs_function_trigger.test.id
}

data "huaweicloud_fgs_function_triggers" "filter_by_trigger_id" {
  function_urn = huaweicloud_fgs_function.test.urn
  trigger_id   = local.trigger_id
}

data "huaweicloud_fgs_function_triggers" "filter_by_not_found_trigger_id" {
  # Query triggers using a not exist trigger ID after trigger resource create.
  depends_on = [
    huaweicloud_fgs_function_trigger.test,
  ]

  function_urn = huaweicloud_fgs_function.test.urn
  trigger_id   = "%[2]s"
}

locals {
  trigger_id_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_trigger_id.triggers[*].id : v == local.trigger_id
  ]
}

output "is_trigger_id_filter_useful" {
  value = length(local.trigger_id_filter_result) > 0 && alltrue(local.trigger_id_filter_result)
}

output "trigger_id_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_function_triggers.filter_by_not_found_trigger_id.triggers) == 0
}

# Filter by trigger type.
locals {
  trigger_type = huaweicloud_fgs_function_trigger.test.type
}

data "huaweicloud_fgs_function_triggers" "filter_by_type" {
  # The behavior of parameter 'type' of the trigger resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_fgs_function_trigger.test,
  ]

  function_urn = huaweicloud_fgs_function.test.urn
  type         = local.trigger_type
}

data "huaweicloud_fgs_function_triggers" "filter_by_not_found_type" {
  # Query triggers using a not exist trigger type after trigger resource create.
  depends_on = [
    huaweicloud_fgs_function_trigger.test,
  ]

  function_urn = huaweicloud_fgs_function.test.urn
  type         = "type_not_found"
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_type.triggers[*].type : v == local.trigger_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

output "type_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_function_triggers.filter_by_not_found_type.triggers) == 0
}

# Filter by trigger status.
locals {
  trigger_status = huaweicloud_fgs_function_trigger.test.status
}

data "huaweicloud_fgs_function_triggers" "filter_by_status" {
  function_urn = huaweicloud_fgs_function.test.urn
  status       = local.trigger_status
}

data "huaweicloud_fgs_function_triggers" "filter_by_not_found_status" {
  function_urn = huaweicloud_fgs_function.test.urn
  status       = "status_not_found"
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_status.triggers[*].status : v == local.trigger_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

output "status_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_function_triggers.filter_by_not_found_status.triggers) == 0
}

# Filter by start time
locals {
  start_time = huaweicloud_fgs_function_trigger.test.created_at
}

data "huaweicloud_fgs_function_triggers" "filter_by_start_time" {
  function_urn = huaweicloud_fgs_function.test.urn
  start_time   = timeadd(local.start_time, "-1s")
}

locals {
  start_time_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_start_time.triggers[*].created_at : timecmp(v, local.start_time) >= 0
  ]
}

output "is_start_time_filter_useful" {
  value = length(local.start_time_filter_result) > 0 && alltrue(local.start_time_filter_result)
}

# Filter by end time
locals {
  end_time = huaweicloud_fgs_function_trigger.test.created_at
}

data "huaweicloud_fgs_function_triggers" "filter_by_end_time" {
  function_urn = huaweicloud_fgs_function.test.urn
  end_time     = timeadd(local.end_time, "1s")
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_end_time.triggers[*].created_at : timecmp(v, local.end_time) <= 0
  ]
}

output "is_end_time_filter_useful" {
  value = length(local.end_time_filter_result) > 0 && alltrue(local.end_time_filter_result)
}

# Specifies the trigger type, query using two modes, and assert that the types of the two data source query results
# are consistent and both match the expected type.
data "huaweicloud_fgs_function_triggers" "filter_by_type_without_function_urn" {
  type = local.trigger_type

  depends_on = [huaweicloud_fgs_function_trigger.test]
}

output "trigger_type_consistency_check_pass" {
  value = alltrue([for v in data.huaweicloud_fgs_function_triggers.filter_by_type.triggers :
    alltrue([for item in data.huaweicloud_fgs_function_triggers.filter_by_type_without_function_urn.triggers :
      item.type == v.type && item.type == local.trigger_type
    ])
  ])
}
`, name, randUUID)
}
