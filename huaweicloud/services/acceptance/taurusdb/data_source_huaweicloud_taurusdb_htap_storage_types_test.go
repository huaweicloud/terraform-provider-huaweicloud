package taurusdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStorageTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_storage_types.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapStorageTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.#"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.0.az_status.%"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.0.min_volume_size"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapStorageTypes_basic() string {
	return `
data "huaweicloud_taurusdb_htap_storage_types" "test" {
  database     = "star-rocks"
  version_name = "3.1.6.0"
}
`
}
