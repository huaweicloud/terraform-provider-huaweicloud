package ga

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessLogs_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_access_logs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byLogId   = "data.huaweicloud_ga_access_logs.filter_by_log_id"
		dcByLogId = acceptance.InitDataSourceCheck(byLogId)

		byStatus   = "data.huaweicloud_ga_access_logs.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byResourceType   = "data.huaweicloud_ga_access_logs.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byResourceIds   = "data.huaweicloud_ga_access_logs.filter_by_resource_ids"
		dcByResourceIds = acceptance.InitDataSourceCheck(byResourceIds)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccessLogs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.log_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.log_stream_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.status"),
					resource.TestMatchResourceAttr(dataSourceName, "logs.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSourceName, "logs.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					dcByLogId.CheckResourceExists(),
					resource.TestCheckOutput("log_id_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("resource_type_filter_useful", "true"),

					dcByResourceIds.CheckResourceExists(),
					resource.TestCheckOutput("resource_ids_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAccessLogs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_access_logs" "test" {
  depends_on = [huaweicloud_ga_access_log.test]
}

locals {
  log_id = data.huaweicloud_ga_access_logs.test.logs[0].id
}

data "huaweicloud_ga_access_logs" "filter_by_log_id" {
  log_id = local.log_id
}

output "log_id_filter_useful" {
  value = length(data.huaweicloud_ga_access_logs.filter_by_log_id.logs) > 0 && alltrue(
    [for v in data.huaweicloud_ga_access_logs.filter_by_log_id.logs[*].id : v == local.log_id]
  )
}

locals {
  status = data.huaweicloud_ga_access_logs.test.logs[0].status
}

data "huaweicloud_ga_access_logs" "filter_by_status" {
  status = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_ga_access_logs.filter_by_status.logs) > 0 && alltrue(
    [for v in data.huaweicloud_ga_access_logs.filter_by_status.logs[*].status : v == local.status]
  )
}

locals {
  resource_type = data.huaweicloud_ga_access_logs.test.logs[0].resource_type
}

data "huaweicloud_ga_access_logs" "filter_by_resource_type" {
  resource_type = local.resource_type
}

output "resource_type_filter_useful" {
  value = length(data.huaweicloud_ga_access_logs.filter_by_resource_type.logs) > 0 && alltrue(
    [for v in data.huaweicloud_ga_access_logs.filter_by_resource_type.logs[*].resource_type : v == local.resource_type]
  )
}

locals {
  resource_ids = [data.huaweicloud_ga_access_logs.test.logs[0].resource_id]
}

data "huaweicloud_ga_access_logs" "filter_by_resource_ids" {
  resource_ids = local.resource_ids
}

output "resource_ids_filter_useful" {
  value = length(data.huaweicloud_ga_access_logs.filter_by_resource_ids.logs) == 1
}
`, testAccessLog_basic(name))
}
