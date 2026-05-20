package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsObjectCompare_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_object_compare.test"
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
				Config: testAccDataSourceDrsObjectCompare_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "export_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "report_remain_seconds"),
					resource.TestCheckResourceAttrSet(dataSourceName, "compare_job_id"),
				),
			},
		},
	})
}

func testAccDataSourceDrsObjectCompare_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_object_compare" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
