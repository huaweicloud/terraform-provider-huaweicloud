package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussSqlTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_sql_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussLimitSqlModels_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "node_limit_sql_model_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_limit_sql_model_list.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "node_limit_sql_model_list.0.sql_model"),
					resource.TestCheckOutput("sql_model_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussLimitSqlModels_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_instances" "test" {}

locals {
  instance       = [for v in  data.huaweicloud_gaussdb_opengauss_instances.test.instances : v if v.id == "%[1]s"][0]
  master_node_id = [for v in local.instance.nodes : v if v.role == "master"][0].id
}

data "huaweicloud_gaussdb_opengauss_sql_templates" "test" {
  instance_id = "%[1]s"
  node_id     = local.master_node_id
}

locals {
  sql_model = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[0].sql_model
}

data "huaweicloud_gaussdb_opengauss_sql_templates" "sql_model_filter" {
  instance_id = "%[1]s"
  node_id     = local.master_node_id
  sql_model   = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[0].sql_model
}

output "sql_model_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_templates.sql_model_filter.node_limit_sql_model_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_sql_templates.sql_model_filter.node_limit_sql_model_list[*].sql_model :
  v == local.sql_model]
  )
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID)
}
