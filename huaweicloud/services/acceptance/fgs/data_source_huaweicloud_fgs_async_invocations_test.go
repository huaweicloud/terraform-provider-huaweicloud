package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please make sure the function includes at least one failed async invocation and one
// success async invocation.
func TestAccDataAsyncInvocations_basic(t *testing.T) {
	var (
		all                 = "data.huaweicloud_fgs_async_invocations.all"
		dcForAllInvocations = acceptance.InitDataSourceCheck(all)

		byRequestId           = "data.huaweicloud_fgs_async_invocations.filter_by_request_id"
		dcByRequestId         = acceptance.InitDataSourceCheck(byRequestId)
		byNotFoundRequestId   = "data.huaweicloud_fgs_async_invocations.filter_by_not_found_request_id"
		dcByNotFoundRequestId = acceptance.InitDataSourceCheck(byNotFoundRequestId)

		bySuccessStatus   = "data.huaweicloud_fgs_async_invocations.filter_by_success_status"
		dcBySuccessStatus = acceptance.InitDataSourceCheck(bySuccessStatus)
		byFailStatus      = "data.huaweicloud_fgs_async_invocations.filter_by_fail_status"
		dcByFailStatus    = acceptance.InitDataSourceCheck(byFailStatus)

		byBeginTime   = "data.huaweicloud_fgs_async_invocations.filter_by_begin_time"
		dcByBeginTime = acceptance.InitDataSourceCheck(byBeginTime)
		byEndTime     = "data.huaweicloud_fgs_async_invocations.filter_by_end_time"
		dcByEndTime   = acceptance.InitDataSourceCheck(byEndTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsFunctionName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAsyncInvocations_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllInvocations.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "invocations.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					// Filter by request ID.
					dcByRequestId.CheckResourceExists(),
					resource.TestCheckOutput("is_request_id_filter_useful", "true"),
					dcByNotFoundRequestId.CheckResourceExists(),
					resource.TestCheckOutput("request_id_not_found_validation_pass", "true"),
					// Filter by status.
					dcBySuccessStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_success_status_filter_useful", "true"),
					dcByFailStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_fail_status_filter_useful", "true"),
					// Filter by time range.
					dcByBeginTime.CheckResourceExists(),
					resource.TestCheckOutput("is_begin_time_filter_useful", "true"),
					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrSet(byRequestId, "invocations.0.request_id"),
					resource.TestCheckResourceAttrSet(byRequestId, "invocations.0.status"),
					resource.TestCheckResourceAttrSet(byRequestId, "invocations.0.start_time"),
					resource.TestCheckResourceAttrSet(byRequestId, "invocations.0.end_time"),
					resource.TestCheckResourceAttrSet(byFailStatus, "invocations.0.error_code"),
					resource.TestCheckResourceAttrSet(byFailStatus, "invocations.0.error_message"),
				),
			},
		},
	})
}

func testAccDataAsyncInvocations_basic() string {
	randRequestId, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_fgs_functions" "test" {
  name = "%[1]s"
}

locals {
  function_urn = try(data.huaweicloud_fgs_functions.test.functions[0].urn, "NOT_FOUND")
}

data "huaweicloud_fgs_async_invocations" "all" {
  function_urn = local.function_urn
}

locals {
  request_id = try(data.huaweicloud_fgs_async_invocations.all.invocations[0].request_id, "NOT_FOUND")
}

# Filter by request ID
data "huaweicloud_fgs_async_invocations" "filter_by_request_id" {
  depends_on = [
    data.huaweicloud_fgs_async_invocations.all,
  ]

  function_urn = local.function_urn
  request_id   = local.request_id
}

locals {
  request_id_filter_result = [
    for v in data.huaweicloud_fgs_async_invocations.filter_by_request_id.invocations[*].request_id : v == local.request_id
  ]
}

output "is_request_id_filter_useful" {
  value = length(local.request_id_filter_result) > 0 && alltrue(local.request_id_filter_result)
}

# Filter by not found request ID
data "huaweicloud_fgs_async_invocations" "filter_by_not_found_request_id" {
  depends_on = [
    data.huaweicloud_fgs_async_invocations.all,
  ]

  function_urn = local.function_urn
  request_id   = "%[2]s"
}

output "request_id_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_async_invocations.filter_by_not_found_request_id.invocations) == 0
}

# Filter by status
locals {
  success_status = "SUCCESS"
  fail_status    = "FAIL"
}

data "huaweicloud_fgs_async_invocations" "filter_by_success_status" {
  depends_on = [
    data.huaweicloud_fgs_async_invocations.all,
  ]

  function_urn = local.function_urn
  status       = local.success_status
}

locals {
  success_status_filter_result = [
    for v in data.huaweicloud_fgs_async_invocations.filter_by_success_status.invocations : alltrue([
      v.status == local.success_status,
      v.error_code == "",
      v.error_message == "",
    ])
  ]
}

output "is_success_status_filter_useful" {
  value = length(local.success_status_filter_result) > 0 && alltrue(local.success_status_filter_result)
}

data "huaweicloud_fgs_async_invocations" "filter_by_fail_status" {
  depends_on = [
    data.huaweicloud_fgs_async_invocations.all,
  ]

  function_urn = local.function_urn
  status       = local.fail_status
}

locals {
  fail_status_filter_result = [
    for v in data.huaweicloud_fgs_async_invocations.filter_by_fail_status.invocations : alltrue([
      v.status == local.fail_status,
      v.error_code != "",
      v.error_message != "",
    ])
  ]
}

output "is_fail_status_filter_useful" {
  value = length(local.fail_status_filter_result) > 0 && alltrue(local.fail_status_filter_result)
}

# Filter by begin time
locals {
  begin_time = try(data.huaweicloud_fgs_async_invocations.all.invocations[0].start_time, "NOT_FOUND")
}

data "huaweicloud_fgs_async_invocations" "filter_by_begin_time" {
  depends_on = [
    data.huaweicloud_fgs_async_invocations.all,
  ]

  function_urn     = local.function_urn
  query_begin_time = local.begin_time
}

locals {
  begin_time_filter_result = [
    for v in data.huaweicloud_fgs_async_invocations.filter_by_begin_time.invocations : timecmp(v.start_time, local.begin_time) <= 0
  ]
}

output "is_begin_time_filter_useful" {
  value = length(local.begin_time_filter_result) > 0 && alltrue(local.begin_time_filter_result)
}

# Filter by end time
locals {
  end_time = try(data.huaweicloud_fgs_async_invocations.all.invocations[0].end_time, "NOT_FOUND")
}

data "huaweicloud_fgs_async_invocations" "filter_by_end_time" {
  depends_on = [
    data.huaweicloud_fgs_async_invocations.all,
  ]

  function_urn   = local.function_urn
  query_end_time = local.end_time
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_fgs_async_invocations.filter_by_end_time.invocations : timecmp(v.end_time, local.end_time) >= 0
  ]
}

output "is_end_time_filter_useful" {
  value = length(local.end_time_filter_result) > 0 && alltrue(local.end_time_filter_result)
}
`, acceptance.HW_FGS_FUNCTION_NAME, randRequestId)
}
