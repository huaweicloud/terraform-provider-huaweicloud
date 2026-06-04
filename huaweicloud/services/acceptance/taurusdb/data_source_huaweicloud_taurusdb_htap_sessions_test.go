package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_sessions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBHtapSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.user"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.command"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.sql_statement"),
					resource.TestCheckResourceAttrSet(dataSource, "process_list.0.duration"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBHtapSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_sessions" "test" {
  instance_id = "%s"
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
