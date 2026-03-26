package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashConfigurations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_logstash_configurations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a logstash cluster and
			// the cluster has configuration files.
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssLogstashConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.conf_content"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.setting.#"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.setting.0.workers"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.setting.0.queue_type"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssLogstashConfigurations_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_logstash_configurations" "test" {
  cluster_id = "%[1]s"
}

locals {
  name   = data.huaweicloud_css_logstash_configurations.test.confs[0].name
  status = data.huaweicloud_css_logstash_configurations.test.confs[0].status
}

data "huaweicloud_css_logstash_configurations" "filter_by_name" {
  cluster_id = "%[1]s"
  name       = local.name
}

data "huaweicloud_css_logstash_configurations" "filter_by_status" {
  cluster_id = "%[1]s"
  status     = local.status
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_css_logstash_configurations.filter_by_name.confs) > 0 && alltrue(
    [for v in data.huaweicloud_css_logstash_configurations.filter_by_name.confs[*].name : v == local.name]
  )
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_css_logstash_configurations.filter_by_status.confs) > 0 && alltrue(
    [for v in data.huaweicloud_css_logstash_configurations.filter_by_status.confs[*].status : v == local.status]
  )
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
