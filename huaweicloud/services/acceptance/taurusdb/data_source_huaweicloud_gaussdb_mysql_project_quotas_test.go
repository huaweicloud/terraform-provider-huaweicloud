package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlProjectQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_project_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlProjectQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.quota"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlProjectQuotas_basic() string {
	return `
data "huaweicloud_gaussdb_mysql_project_quotas" "test" {}

locals {
  type = "instance"
}
data "huaweicloud_gaussdb_mysql_project_quotas" "type_filter" {
  type = "instance"
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_project_quotas.type_filter.quotas) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_project_quotas.type_filter.quotas[*].resources : length(v) > 0 && alltrue(
  [for vv in v : vv.type == local.type]
  )]
  )
}
`
}
