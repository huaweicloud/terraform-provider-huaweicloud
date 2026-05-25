package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsDataProcessingRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_data_processing_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsDataProcessingRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_process_info.#"),
				),
			},
		},
	})
}

func testAccDataSourceDrsDataProcessingRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_data_processing_rules" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_IDS)
}
