package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashPipelines_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_logstash_pipelines.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCssLogstashPipelines_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pipelines.0.name"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCssLogstashPipelines_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_logstash_pipelines" "test" {
  depends_on = [huaweicloud_css_logstash_pipeline.test]

  cluster_id = huaweicloud_css_logstash_cluster.test.id
}

locals {
  name = data.huaweicloud_css_logstash_pipelines.test.pipelines[0].name
}

data "huaweicloud_css_logstash_pipelines" "filter_by_name" {
  cluster_id =  huaweicloud_css_logstash_cluster.test.id
  name       = local.name
}

locals {
  list_by_name = data.huaweicloud_css_logstash_pipelines.filter_by_name.pipelines
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}
`, testLogstashPipeline_basic(name))
}
