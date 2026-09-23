package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsExtraColumnInfo_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_extra_column_info.test"
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
				Config: testAccDataSourceDrsExtraColumnInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "column_process_objects.#"),
				),
			},
		},
	})
}

func testAccDataSourceDrsExtraColumnInfo_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_extra_column_info" "test" {
  job_id = "%[1]s"
}
`, acceptance.HW_DRS_JOB_ID)
}
