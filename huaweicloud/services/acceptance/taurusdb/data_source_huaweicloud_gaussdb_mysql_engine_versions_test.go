package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlEngineVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_engine_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlEngineVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.kernel_version"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlEngineVersions_basic() string {
	return `
data "huaweicloud_gaussdb_mysql_engine_versions" "test" {
  database_name = "gaussdb-mysql"
}
`
}
