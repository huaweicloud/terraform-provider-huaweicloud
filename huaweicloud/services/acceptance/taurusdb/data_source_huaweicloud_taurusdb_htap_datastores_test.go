package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapDatastores_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_datastores.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapDatastores_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "datastores.0.kernel_version"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapDatastores_basic() string {
	return `
data "huaweicloud_taurusdb_htap_datastores" "test" {
  engine_name = "star-rocks"
}
`
}
