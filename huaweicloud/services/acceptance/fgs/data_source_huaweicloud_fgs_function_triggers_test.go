package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFunctionTriggers_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_fgs_function_triggers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byId   = "data.huaweicloud_fgs_function_triggers.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byType   = "data.huaweicloud_fgs_function_triggers.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus   = "data.huaweicloud_fgs_function_triggers.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byCreatedAt   = "data.huaweicloud_fgs_function_triggers.filter_by_created_at"
		dcByCreatedAt = acceptance.InitDataSourceCheck(byCreatedAt)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTriggers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "triggers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestMatchResourceAttr(dataSource, "triggers.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "triggers.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByCreatedAt.CheckResourceExists(),
					resource.TestCheckOutput("is_created_at_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccFunctionTriggers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_fgs_function_triggers" "test" {
  depends_on = [
    huaweicloud_fgs_function_trigger.timer_rate
  ]

  function_urn = huaweicloud_fgs_function.test.urn
}

// Filter by ID
locals {
  trigger_id = huaweicloud_fgs_function_trigger.timer_rate.id
}

data "huaweicloud_fgs_function_triggers" "filter_by_id" {
  function_urn = huaweicloud_fgs_function.test.urn
  trigger_id   = local.trigger_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_id.triggers[*].id : v == local.trigger_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

// Filter by type
locals {
  trigger_type = data.huaweicloud_fgs_function_triggers.test.triggers[0].type
}

data "huaweicloud_fgs_function_triggers" "filter_by_type" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = local.trigger_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_type.triggers[*].type : v == local.trigger_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

// Filter by status
locals {
  trigger_status = huaweicloud_fgs_function_trigger.timer_rate.status
}

data "huaweicloud_fgs_function_triggers" "filter_by_status" {
  function_urn = huaweicloud_fgs_function.test.urn
  status       = local.trigger_status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_status.triggers[*].status : v == local.trigger_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

// Filter by creation time
locals {
  created_time = huaweicloud_fgs_function_trigger.timer_rate.created_at
}

data "huaweicloud_fgs_function_triggers" "filter_by_created_at" {
  function_urn = huaweicloud_fgs_function.test.urn
  start_time   = local.created_time
  end_time     = local.created_time
}

locals {
  created_at_filter_result = [
    for v in data.huaweicloud_fgs_function_triggers.filter_by_created_at.triggers[*].created_at : timecmp(v, local.created_time) == 0
  ]
}

output "is_created_at_filter_useful" {
  value = length(local.created_at_filter_result) > 0 && alltrue(local.created_at_filter_result)
}
`, testAccFunctionTimingTrigger_basic_step1(name))
}
