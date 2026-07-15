package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAuthorizeDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_authorize_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscAuthorizeDatabases_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.asset_name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.auth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.authorized"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_address"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_authorized"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_port"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_user"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_version"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.default"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.ins_id"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.ins_type"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.region"),
				),
			},
		},
	})
}

const testAccDataSourceDscAuthorizeDatabases_basic = `
data "huaweicloud_dsc_authorize_databases" "test" {
  instance_type = "RDS"
}
`
