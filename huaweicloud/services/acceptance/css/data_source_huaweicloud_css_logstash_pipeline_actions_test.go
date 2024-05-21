package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashPipelineActions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_logstash_pipeline_actions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssLogstashPipelineActions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "actions.0.status"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssLogstashPipelineActions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_logstash_pipeline_actions" "test" {
  depends_on = [huaweicloud_css_logstash_pipeline.test]

  cluster_id = huaweicloud_css_logstash_cluster.test.id
}

locals {
  action_id = data.huaweicloud_css_logstash_pipeline_actions.test.actions[0].id
  type      = data.huaweicloud_css_logstash_pipeline_actions.test.actions[0].type
  status    = data.huaweicloud_css_logstash_pipeline_actions.test.actions[0].status
}

data "huaweicloud_css_logstash_pipeline_actions" "filter_by_id" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  action_id  = local.action_id
}

data "huaweicloud_css_logstash_pipeline_actions" "filter_by_type" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  type       = local.type
}

data "huaweicloud_css_logstash_pipeline_actions" "filter_by_status" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  status     = local.status
}

locals {
  list_by_id     = data.huaweicloud_css_logstash_pipeline_actions.filter_by_id.actions
  list_by_type   = data.huaweicloud_css_logstash_pipeline_actions.filter_by_type.actions
  list_by_status = data.huaweicloud_css_logstash_pipeline_actions.filter_by_status.actions
}

output "id_filter_is_useful" {
  value = length(local.list_by_id) > 0 && alltrue(
    [for v in local.list_by_id[*].id : v == local.action_id]
  )
}

output "type_filter_is_useful" {
  value = length(local.list_by_type) > 0 && alltrue(
    [for v in local.list_by_type[*].type : v == local.type]
  )
}

output "status_filter_is_useful" {
  value = length(local.list_by_type) > 0 && alltrue(
    [for v in local.list_by_type[*].status : v == local.status]
  )
}
`, testLogstashPipeline_basic(name))
}
