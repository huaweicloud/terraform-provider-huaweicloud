package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbMysqlConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.user_defined"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbMysqlConfigurations_basic() string {
	return `
data "huaweicloud_gaussdb_mysql_configurations" "test" {}
`
}
