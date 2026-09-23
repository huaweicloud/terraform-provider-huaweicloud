package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCompareUsers_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_compare_users.test"
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
				Config: testAccDataSourceCompareUsers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_compare_info.#"),
				),
			},
		},
	})
}

func testAccDataSourceCompareUsers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_users" "test" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}
