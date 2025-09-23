package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOperationLogs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_cbr_operation_logs.basic"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOperationLogs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "operation_logs.0.operation_type", "vault_delete"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.ended_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.started_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.provider_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.vault_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.vault_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.extra_info.0.vault_delete.0.fail_delete_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.extra_info.0.vault_delete.0.total_delete_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.extra_info.0.common.0.progress"),
					resource.TestCheckResourceAttrSet(dataSourceName, "operation_logs.0.extra_info.0.common.0.request_id"),
					resource.TestCheckOutput("is_operation_type_filter_useful", "true"),
					resource.TestCheckOutput("is_provider_id_filter_useful", "true"),
					resource.TestCheckOutput("is_vault_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_vault_name_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceOperationLogs_basic = `
data "huaweicloud_cbr_operation_logs" "basic" {
  status         = "success"
  operation_type = "vault_delete"
}

locals {
  filtered_logs = [
    for log in data.huaweicloud_cbr_operation_logs.basic.operation_logs : log
    if log.operation_type == "vault_delete" && 
       length(log.extra_info) > 0 && 
       log.extra_info[0].vault_delete[0].total_delete_count == 1
  ]
  
  operation_type = length(local.filtered_logs) > 0 ? 
    local.filtered_logs[0].operation_type : 
    data.huaweicloud_cbr_operation_logs.basic.operation_logs[0].operation_type
  provider_id = length(local.filtered_logs) > 0 ? 
    local.filtered_logs[0].provider_id : 
    data.huaweicloud_cbr_operation_logs.basic.operation_logs[0].provider_id
  vault_id = length(local.filtered_logs) > 0 ? 
    local.filtered_logs[0].vault_id : 
    data.huaweicloud_cbr_operation_logs.basic.operation_logs[0].vault_id
  status = length(local.filtered_logs) > 0 ? 
    local.filtered_logs[0].status : 
    data.huaweicloud_cbr_operation_logs.basic.operation_logs[0].status
  vault_name = length(local.filtered_logs) > 0 ? 
    local.filtered_logs[0].vault_name : 
    data.huaweicloud_cbr_operation_logs.basic.operation_logs[0].vault_name
}

data "huaweicloud_cbr_operation_logs" "filter_by_operation_type" {
  operation_type = local.operation_type
}

output "is_operation_type_filter_useful" {
  value = length(data.huaweicloud_cbr_operation_logs.filter_by_operation_type.operation_logs) > 0 && alltrue([
    for v in data.huaweicloud_cbr_operation_logs.filter_by_operation_type.operation_logs[*].operation_type : v == local.operation_type
  ])
}

data "huaweicloud_cbr_operation_logs" "filter_by_provider_id" {
  provider_id = local.provider_id
}

output "is_provider_id_filter_useful" {
  value = length(data.huaweicloud_cbr_operation_logs.filter_by_provider_id.operation_logs) > 0 && alltrue([
    for v in data.huaweicloud_cbr_operation_logs.filter_by_provider_id.operation_logs[*].provider_id : v == local.provider_id
  ])
}

data "huaweicloud_cbr_operation_logs" "filter_by_vault_id" {
  vault_id = local.vault_id
}

output "is_vault_id_filter_useful" {
  value = length(data.huaweicloud_cbr_operation_logs.filter_by_vault_id.operation_logs) > 0 && alltrue([
    for v in data.huaweicloud_cbr_operation_logs.filter_by_vault_id.operation_logs[*].vault_id : v == local.vault_id
  ])
}

data "huaweicloud_cbr_operation_logs" "filter_by_status" {
  status   = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_cbr_operation_logs.filter_by_status.operation_logs) > 0 && alltrue([
    for v in data.huaweicloud_cbr_operation_logs.filter_by_status.operation_logs[*].status : v == local.status
  ])
}

data "huaweicloud_cbr_operation_logs" "filter_by_vault_name" {
  vault_name = local.vault_name
}

output "is_vault_name_filter_useful" {
  value = length(data.huaweicloud_cbr_operation_logs.filter_by_vault_name.operation_logs) > 0 ? "true" : "false"
}
`
