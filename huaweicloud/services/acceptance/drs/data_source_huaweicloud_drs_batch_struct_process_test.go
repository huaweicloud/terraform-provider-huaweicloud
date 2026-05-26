package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsBatchStructProcess_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_batch_struct_process.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDrsBatchStructProcess_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_process.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_process.0.result.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_process.0.result.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_process.0.result.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_process.0.create_time"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDrsBatchStructProcess_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_batch_struct_process" "test" {
  job_ids = split(",", "%s")
}
`, acceptance.HW_DRS_JOB_IDS)
}
