package ddm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmConfigurations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.updated"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.user_defined"),
				),
			},
		},
	})
}

func testAccDatasourceDdmConfigurations_basic() string {
	return `
data "huaweicloud_ddm_configurations" "test" {}
`
}
