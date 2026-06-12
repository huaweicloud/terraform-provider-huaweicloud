package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksReplications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_replications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksReplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "replications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "replications.0.source_database"),
					resource.TestCheckResourceAttrSet(dataSource, "replications.0.target_database"),
					resource.TestCheckResourceAttrSet(dataSource, "replications.0.task_name"),
					resource.TestCheckResourceAttr(dataSource, "replications.0.status", "Yes"),
					resource.TestCheckResourceAttr(dataSource, "replications.0.stage", "Wait"),
					resource.TestCheckResourceAttr(dataSource, "replications.0.percentage", "0"),
					resource.TestCheckResourceAttr(dataSource, "replications.0.is_need_repair", "false"),
					resource.TestCheckResourceAttr(dataSource, "replications.0.is_main_task", "true"),
				),
			},
		},
	})
}

func testAccDataSourceTaurusDBHtapStarrocksReplications_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_replications" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID)
}
