package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsBatchStructDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_batch_struct_detail.test"
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
				Config: testAccDataSourceDrsBatchStructDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_detail.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_detail.0.total_record"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.struct_detail.0.list.0.progress"),
				),
			},
		},
	})
}

func testAccDataSourceDrsBatchStructDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_batch_struct_detail" "test" {
  type = "database"
  jobs = split(",", "%s")
}
`, acceptance.HW_DRS_JOB_IDS)
}
