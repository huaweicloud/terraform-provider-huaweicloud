package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashPipelineActions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_logstash_pipeline_actions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a logstash cluster and
			// there are active pipelines in this cluster.
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssLogstashPipelineActions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.updated_at"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssLogstashPipelineActions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_logstash_pipeline_actions" "test" {
  cluster_id = "%[1]s"
}

locals {
  action_id = data.huaweicloud_css_logstash_pipeline_actions.test.actions[0].id
  type      = data.huaweicloud_css_logstash_pipeline_actions.test.actions[0].type
  status    = data.huaweicloud_css_logstash_pipeline_actions.test.actions[0].status
}

data "huaweicloud_css_logstash_pipeline_actions" "filter_by_id" {
  cluster_id = "%[1]s"
  action_id  = local.action_id
}

data "huaweicloud_css_logstash_pipeline_actions" "filter_by_type" {
  cluster_id = "%[1]s"
  type       = local.type
}

data "huaweicloud_css_logstash_pipeline_actions" "filter_by_status" {
  cluster_id = "%[1]s"
  status     = local.status
}

output "id_filter_is_useful" {
  value = length(data.huaweicloud_css_logstash_pipeline_actions.filter_by_id.actions) > 0 && alltrue(
    [for v in data.huaweicloud_css_logstash_pipeline_actions.filter_by_id.actions[*].id : v == local.action_id]
  )
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_css_logstash_pipeline_actions.filter_by_type.actions) > 0 && alltrue(
    [for v in data.huaweicloud_css_logstash_pipeline_actions.filter_by_type.actions[*].type : v == local.type]
  )
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_css_logstash_pipeline_actions.filter_by_status.actions) > 0 && alltrue(
    [for v in data.huaweicloud_css_logstash_pipeline_actions.filter_by_status.actions[*].status : v == local.status]
  )
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
