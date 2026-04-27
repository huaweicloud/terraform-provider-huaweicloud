package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPublicationMonitor_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_publication_monitor.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsPublicationMonitor_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "worst_latency"),
					resource.TestCheckResourceAttrSet(dataSource, "best_latency"),
					resource.TestCheckResourceAttrSet(dataSource, "average_latency"),
					resource.TestCheckResourceAttrSet(dataSource, "last_dist_sync"),
					resource.TestCheckResourceAttrSet(dataSource, "replicated_transactions"),
					resource.TestCheckResourceAttrSet(dataSource, "replication_rate_trans"),
				),
			},
		},
	})
}

func testAccDataSourceRdsPublicationMonitor_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_publications" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_publication_monitor" "test" {
  instance_id    = "%[1]s"
  publication_id = data.huaweicloud_rds_publications.test.publications[0].id
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
