package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsCompareProgress_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_compare_progress.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
			acceptance.TestAccPreCheckDrsCompareJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsCompareProgress_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_info.0.progress"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_info.0.src_speed"),
					resource.TestCheckResourceAttrSet(dataSourceName, "full_info.0.recheck_entities"),
					resource.TestCheckResourceAttrSet(dataSourceName, "incre_info.0.delay"),
					resource.TestCheckResourceAttrSet(dataSourceName, "incre_info.0.src_speed"),
					resource.TestCheckResourceAttrSet(dataSourceName, "incre_info.0.rps"),
					resource.TestCheckResourceAttrSet(dataSourceName, "incre_info.0.recheck_entities"),
					resource.TestCheckResourceAttrSet(dataSourceName, "global_info.0.src_speed"),
				),
			},
		},
	})
}

func testAccDataSourceDrsCompareProgress_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_progress" "test" {
  job_id         = "%s"
  compare_job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}
