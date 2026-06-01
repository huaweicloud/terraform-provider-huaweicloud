package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsBatchRposAndRtos_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_batch_rpos_and_rtos.test"
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
				Config: testDataSourceDrsBatchRposAndRtos_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rpo_info.0.check_point"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rpo_info.0.delay"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rpo_info.0.gtid_set"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rpo_info.0.time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rto_info.0.check_point"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rto_info.0.delay"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rto_info.0.gtid_set"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.rto_info.0.time"),
				),
			},
		},
	})
}

func testDataSourceDrsBatchRposAndRtos_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_batch_rpos_and_rtos" "test" {
  jobs = split(",", "%s")
}
`, acceptance.HW_DRS_JOB_IDS)
}
