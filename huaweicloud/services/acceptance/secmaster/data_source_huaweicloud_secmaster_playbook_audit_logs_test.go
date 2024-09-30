package secmaster

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccDataSourceSecmasterPlaybookAuditLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbook_audit_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybookAuditLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.action_id"),

					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_action_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybookAuditLogs_basic() string {
	nowStamp := time.Now().Unix()
	startTime := utils.FormatTimeStampRFC3339(nowStamp-60*60, false, "2006-01-02T15:04:05.000Z+0800")
	endTime := utils.FormatTimeStampRFC3339(nowStamp, false, "2006-01-02T15:04:05.000Z+0800")
	return fmt.Sprintf(`
data "huaweicloud_secmaster_playbook_audit_logs" "test" {
  workspace_id = "%[1]s"
  start_time   = "%[2]s"
  end_time     = "%[3]s"
}

locals {
  instance_id = data.huaweicloud_secmaster_playbook_audit_logs.test.audit_logs[0].instance_id
  action_id   = data.huaweicloud_secmaster_playbook_audit_logs.test.audit_logs[0].action_id
  status      = data.huaweicloud_secmaster_playbook_audit_logs.test.audit_logs[0].status
}

data "huaweicloud_secmaster_playbook_audit_logs" "filter_by_instance_id" {
  workspace_id = "%[1]s"
  start_time   = "%[2]s"
  end_time     = "%[3]s"
  instance_id  = local.instance_id
}

data "huaweicloud_secmaster_playbook_audit_logs" "filter_by_action_id" {
  workspace_id = "%[1]s"
  start_time   = "%[2]s"
  end_time     = "%[3]s"
  action_id    = local.action_id
}

data "huaweicloud_secmaster_playbook_audit_logs" "filter_by_status" {
  workspace_id = "%[1]s"
  start_time   = "%[2]s"
  end_time     = "%[3]s"
  status       = local.status
}

locals {
  list_by_instance_id = data.huaweicloud_secmaster_playbook_audit_logs.filter_by_instance_id.audit_logs
  list_by_action_id   = data.huaweicloud_secmaster_playbook_audit_logs.filter_by_action_id.audit_logs
  list_by_status      = data.huaweicloud_secmaster_playbook_audit_logs.filter_by_status.audit_logs
}

output "is_instance_id_filter_useful" {
  value = length(local.list_by_instance_id) > 0 && alltrue(
    [for v in local.list_by_instance_id[*].instance_id : v == local.instance_id]
  )
}

output "is_action_id_filter_useful" {
  value = length(local.list_by_action_id) > 0 && alltrue(
    [for v in local.list_by_action_id[*].action_id : v == local.action_id]
  )
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, startTime, endTime)
}
