package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksReplicationDBParameters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_replication_db_parameters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksReplicationDBParameters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "db_parameters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "db_parameters.0.param_name"),
					resource.TestCheckResourceAttrSet(dataSource, "db_parameters.0.data_type"),
					resource.TestCheckResourceAttrSet(dataSource, "db_parameters.0.default_value"),
					resource.TestCheckResourceAttrSet(dataSource, "db_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSource, "db_parameters.0.description"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBHtapStarrocksReplicationDBParameters_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_replication_db_parameters" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
