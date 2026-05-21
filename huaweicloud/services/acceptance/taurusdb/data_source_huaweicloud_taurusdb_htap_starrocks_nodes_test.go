package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksNodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapStarrocksNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "node_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_list.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "node_list.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "node_list.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "node_list.0.status"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapStarrocksNodes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_nodes" "test" {
  instance_id = "%s"
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
