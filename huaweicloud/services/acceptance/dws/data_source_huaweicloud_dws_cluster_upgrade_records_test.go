package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterUpgradeRecords_basic(t *testing.T) {
	all := "data.huaweicloud_dws_cluster_upgrade_records.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataClusterUpgradeRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile(`^[0-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "records.0.id"),
					resource.TestCheckResourceAttrSet(all, "records.0.status"),
					resource.TestCheckResourceAttrSet(all, "records.0.record_type"),
					resource.TestCheckResourceAttrSet(all, "records.0.from_version"),
					resource.TestCheckResourceAttrSet(all, "records.0.to_version"),
					resource.TestCheckResourceAttrSet(all, "records.0.start_time"),
					resource.TestCheckResourceAttrSet(all, "records.0.end_time"),
					resource.TestCheckResourceAttrSet(all, "records.0.job_id"),
					resource.TestCheckResourceAttrSet(all, "records.0.failed_reason"),
				),
			},
		},
	})
}

func testAccDataClusterUpgradeRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_upgrade_records" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
