package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMetricResults_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_metric_results.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMaster(t)
			acceptance.TestAccPreCheckSecMasterMetricID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMetricResults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "metric_results.0.id", acceptance.HW_SECMASTER_METRIC_ID),
					resource.TestCheckResourceAttrSet(dataSource, "metric_results.0.labels.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metric_results.0.data_rows.#"),
				),
			},
		},
	})
}

func testDataSourceMetricResults_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_metric_results" "test" {
  workspace_id = "%[1]s"
  metric_ids   = ["%[2]s"]
  timespan     = "2024-07-12T16:00:00.000Z/2024-08-13T15:59:59.999Z"
  cache        = "true"

  params = [
    {
      start_date = "2024-07-25T00:00:00.000+08:00"
    }
  ]
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_METRIC_ID)
}
