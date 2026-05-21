package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsProgressData_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_progress_data.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDrsProgressData_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_compare_data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_compare_data.0.src_db"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_compare_data.0.src_tb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_compare_data.0.dst_db"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_compare_data.0.dst_tb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flow_compare_data.0.progress"),
				),
			},
		},
	})
}

func testDataSourceDrsProgressData_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_progress_data" "test" {
  job_id = "%s"
  type   = "table"
}
`, acceptance.HW_DRS_JOB_ID)
}
