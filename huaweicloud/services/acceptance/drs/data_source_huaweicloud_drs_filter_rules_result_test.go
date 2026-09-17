package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsFilterRulesResult_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_filter_rules_result.test"
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
				Config: testAccDataSourceDrsFilterRulesResult_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_count"),
				),
			},
		},
	})
}

func testAccDataSourceDrsFilterRulesResult_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_filter_rules_result" "test" {
  job_id = "%[1]s"
}
`, acceptance.HW_DRS_JOB_ID)
}
