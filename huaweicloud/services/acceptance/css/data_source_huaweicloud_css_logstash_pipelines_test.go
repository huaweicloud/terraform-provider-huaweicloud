package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashPipelines_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_logstash_pipelines.test"
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
				Config: testDataSourceCssLogstashPipelines_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.keep_alive"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.update_at"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssLogstashPipelines_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_logstash_pipelines" "test" {
  cluster_id = "%[1]s"
}

locals {
  name = data.huaweicloud_css_logstash_pipelines.test.pipelines[0].name
}

data "huaweicloud_css_logstash_pipelines" "filter_by_name" {
  cluster_id = "%[1]s"
  name       = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_css_logstash_pipelines.filter_by_name.pipelines) > 0 && alltrue(
    [for v in data.huaweicloud_css_logstash_pipelines.filter_by_name.pipelines[*].name : v == local.name]
  )
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
