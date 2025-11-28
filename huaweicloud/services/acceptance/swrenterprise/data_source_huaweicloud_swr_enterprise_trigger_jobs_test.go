package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseTriggerJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_trigger_jobs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseTriggerJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_detail"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.event_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.notify_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseTriggerJobs_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_triggers" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_trigger_jobs" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = data.huaweicloud_swr_enterprise_triggers.test.triggers[0].namespace_name
  policy_id      = data.huaweicloud_swr_enterprise_triggers.test.triggers[0].id
}
`
}
