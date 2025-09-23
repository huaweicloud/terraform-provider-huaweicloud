package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlRecyclingInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_recycling_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlRecyclingInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.ha_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.pay_model"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.deleted_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_vip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.backup_level"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.recycle_backup_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.recycle_status"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlRecyclingInstances_basic() string {
	return `
data "huaweicloud_gaussdb_mysql_recycling_instances" "test" {}
`
}
