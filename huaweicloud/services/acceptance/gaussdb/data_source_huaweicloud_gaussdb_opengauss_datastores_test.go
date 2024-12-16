package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussDatastores_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_datastores.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussDatastores_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.supported_versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.supported_versions.0"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.instance_mode"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussDatastores_basic() string {
	return `
data "huaweicloud_gaussdb_opengauss_datastores" "test" {}
`
}
