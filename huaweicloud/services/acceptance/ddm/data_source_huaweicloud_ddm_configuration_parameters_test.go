package ddm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdmConfigurationParameters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_configuration_parameters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmConfigurationParameters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "created"),
					resource.TestCheckResourceAttrSet(dataSource, "updated"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.restart_required"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.readonly"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSource, "configuration_parameters.0.type"),
				),
			},
		},
	})
}

func testAccDatasourceDdmConfigurationParameters_basic() string {
	return `
data "huaweicloud_ddm_configurations" "test" {}

data "huaweicloud_ddm_configuration_parameters" "test" {
  config_id = data.huaweicloud_ddm_configurations.test.configurations[0].id
}
`
}
