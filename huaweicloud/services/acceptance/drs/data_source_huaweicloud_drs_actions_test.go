package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsJobActions_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_actions.test"
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
				Config: testAccDataSourceDrsJobActions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_action.0.available_actions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "job_action.0.unavailable_actions.#"),
				),
			},
		},
	})
}

func testAccDataSourceDrsJobActions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_actions" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
