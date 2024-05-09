package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_logstash_configurations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssLogstashConfigurations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "confs.0.status"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssLogstashConfigurations_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_logstash_configurations" "test" {
  depends_on = [
    huaweicloud_css_logstash_configuration.test_1,
    huaweicloud_css_logstash_configuration.test_2,
    huaweicloud_css_logstash_configuration.test_3,
  ]

  cluster_id = huaweicloud_css_logstash_cluster.test.id
}

locals {
  name   = data.huaweicloud_css_logstash_configurations.test.confs[0].name
  status = data.huaweicloud_css_logstash_configurations.test.confs[0].status
}

data "huaweicloud_css_logstash_configurations" "filter_by_name" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  name       = local.name
}

data "huaweicloud_css_logstash_configurations" "filter_by_status" {
  cluster_id = huaweicloud_css_logstash_cluster.test.id
  status     = local.status
}

locals {
  list_by_name   = data.huaweicloud_css_logstash_configurations.filter_by_name.confs
  list_by_status = data.huaweicloud_css_logstash_configurations.filter_by_status.confs
}

output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}

output "status_filter_is_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}
`, logstashCluster_configurations(name))
}
