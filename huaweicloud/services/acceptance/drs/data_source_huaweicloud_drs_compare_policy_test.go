package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsComparePolicy_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_compare_policy.test"
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
				Config: testAccDataSourceDrsComparePolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "interval_hour"),
					resource.TestCheckResourceAttrSet(dataSourceName, "period"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "begin_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "end_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "compare_type.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "next_compare_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "compare_policy"),
				),
			},
		},
	})
}

func testAccDataSourceDrsComparePolicy_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_policy" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
