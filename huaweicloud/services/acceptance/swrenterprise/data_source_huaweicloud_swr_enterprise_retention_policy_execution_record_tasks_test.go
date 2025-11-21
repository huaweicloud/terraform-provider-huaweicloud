package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseRetentionPolicyExecutionRecordTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseRetentionPolicyExecutionRecordTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.execution_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status_code"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status_revision"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.retained"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.total"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseRetentionPolicyExecutionRecordTasks_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_enterprise_instances" "test" {}

resource "huaweicloud_swr_enterprise_retention_policy" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  name           = "%[1]s"
  algorithm      = "or"
  enabled        = true
  
  rules {
    priority        = 0
    action          = "retain"
    repo_scope_mode = "regular"
    disabled        = false
    template        = "latestPushedK"

    params = {
      latestPushedK = jsonencode(1)
    }

    scope_selectors {
      key = "repository"

      value {
        kind       = "doublestar"
        decoration = "repoMatches"
        pattern    = "**"
      }
    }
    
    tag_selectors {
      kind       = "doublestar"
      decoration = "matches"
      pattern    = "**"
    }
  }

  trigger {
    type = "scheduled"

    trigger_settings {
      cron = "0 0 0 1 * ?"
    }
  }
}

resource "huaweicloud_swr_enterprise_retention_policy_execute" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_retention_policy.test.policy_id
  dry_run        = true

  provisioner "local-exec" {
    command = "sleep 60"
  }
}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_tasks" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_retention_policy.test.policy_id
  execution_id   = huaweicloud_swr_enterprise_retention_policy_execute.test.execution_id
}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_tasks" "filter_by_status" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_retention_policy.test.policy_id
  execution_id   = huaweicloud_swr_enterprise_retention_policy_execute.test.execution_id
  status         = data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.test.tasks[0].status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.filter_by_status.tasks) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.filter_by_status.tasks[*].status : 
	  v == data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.test.tasks[0].status]
  )
}
`, name)
}
