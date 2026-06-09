package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksParameters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_parameters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapStarrocksParameters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_version_name"),
					resource.TestCheckResourceAttrSet(dataSource, "configurations.0.datastore_name"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.#"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.restart_required"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.readonly"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "parameter_values.0.description"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapStarrocksParameters_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_parameters" "test" {
  instance_id = "%[1]s"
  node_type   = "be"
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
