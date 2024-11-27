package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBRedisFlavorsDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_redis_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbRedisFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.az_status.%"),
				),
			},
		},
	})
}

func testDataSourceGaussdbRedisFlavors_basic() string {
	return `
data "huaweicloud_gaussdb_redis_flavors" "test" {}
`
}
