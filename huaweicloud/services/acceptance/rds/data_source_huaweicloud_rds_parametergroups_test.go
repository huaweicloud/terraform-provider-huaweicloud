package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceParametergroups_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_parametergroups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceParametergroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "configurations.#"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.user_defined"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "configurations.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("datastore_version_name_filter_is_useful", "true"),
					resource.TestCheckOutput("datastore_name_filter_is_useful", "true"),
					resource.TestCheckOutput("user_defined_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceParametergroups_basic() string {
	return `
data "huaweicloud_rds_parametergroups" "test" {
}

data "huaweicloud_rds_parametergroups" "name_filter" {
  name = data.huaweicloud_rds_parametergroups.test.configurations.0.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_parametergroups.test.configurations) > 0 && alltrue(
  [for v in  data.huaweicloud_rds_parametergroups.name_filter.configurations[*].name :
  v == data.huaweicloud_rds_parametergroups.name_filter.name]
)
}

data "huaweicloud_rds_parametergroups" "datastore_version_name_filter" {
  datastore_version_name = data.huaweicloud_rds_parametergroups.test.configurations.0.datastore_version_name
}
output "datastore_version_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_parametergroups.test.configurations) > 0 && alltrue(
  [for v in  data.huaweicloud_rds_parametergroups.datastore_version_name_filter.configurations[*].datastore_version_name :
  v == data.huaweicloud_rds_parametergroups.datastore_version_name_filter.datastore_version_name]
)
}

data "huaweicloud_rds_parametergroups" "datastore_name_filter" {
  datastore_name = data.huaweicloud_rds_parametergroups.test.configurations.0.datastore_name
}
output "datastore_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_parametergroups.test.configurations) > 0 && alltrue(
  [for v in  data.huaweicloud_rds_parametergroups.datastore_name_filter.configurations[*].datastore_name :
  v == data.huaweicloud_rds_parametergroups.datastore_name_filter.datastore_name]
)
}

data "huaweicloud_rds_parametergroups" "user_defined_filter" {
  user_defined = data.huaweicloud_rds_parametergroups.test.configurations.0.user_defined
}
output "user_defined_filter_is_useful" {
  value = length(data.huaweicloud_rds_parametergroups.test.configurations) > 0 && alltrue(
  [for v in  data.huaweicloud_rds_parametergroups.user_defined_filter.configurations[*].user_defined :
  v == data.huaweicloud_rds_parametergroups.user_defined_filter.user_defined]
)
}
`
}
