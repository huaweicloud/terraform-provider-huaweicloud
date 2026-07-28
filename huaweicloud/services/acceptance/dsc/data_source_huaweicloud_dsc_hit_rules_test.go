package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscHitRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_hit_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscScanJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscHitRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hit_rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hit_rules.0.rule_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hit_rules.0.rule_name"),
					resource.TestCheckResourceAttrSet(dataSource, "hit_rules.0.count"),
				),
			},
		},
	})
}

func testAccDataSourceDscHitRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_hit_rules" "test" {
  job_id = "%[1]s"
}
`, acceptance.HW_DSC_SCAN_JOB_ID)
}
