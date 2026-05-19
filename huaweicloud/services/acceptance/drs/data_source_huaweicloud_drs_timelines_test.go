package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsTimelines_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_timelines.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsTimelines_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "timelines.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "timelines.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "timelines.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "timelines.0.operation_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "timelines.0.user_name"),
				),
			},
		},
	})
}

func testAccDataSourceDrsTimelines_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_timelines" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
