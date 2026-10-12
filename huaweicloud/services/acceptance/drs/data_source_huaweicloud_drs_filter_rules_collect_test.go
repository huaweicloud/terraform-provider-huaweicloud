package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsFilterRulesCollect_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_filter_rules_collect.test"
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
				Config: testAccDataSourceDrsFilterRulesCollect_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "job_id", acceptance.HW_DRS_JOB_ID),
					resource.TestCheckResourceAttr(dataSourceName, "filter_type", "common"),
					resource.TestCheckResourceAttrSet(dataSourceName, "query_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_count"),
				),
			},
		},
	})
}

func testAccDataSourceDrsFilterRulesCollect_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_filter_rules_collect" "test" {
  job_id       = "%[1]s"
  filter_type  = "common"
  is_all_rules = true
}
`, acceptance.HW_DRS_JOB_ID)
}
