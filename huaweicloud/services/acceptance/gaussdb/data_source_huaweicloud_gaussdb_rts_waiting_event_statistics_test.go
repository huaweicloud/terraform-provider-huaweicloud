package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbRtsWaitingEventStatistics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_rts_waiting_event_statistics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbRtsWaitingEventStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "wait_events_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "wait_events_info.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "wait_events_info.0.event_name"),
					resource.TestCheckResourceAttrSet(dataSource, "wait_events_info.0.count"),
				),
			},
		},
	})
}

func testDataSourceGaussdbRtsWaitingEventStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_rts_waiting_event_statistics" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
