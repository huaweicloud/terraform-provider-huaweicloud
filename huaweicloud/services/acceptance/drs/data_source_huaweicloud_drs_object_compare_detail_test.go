package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceObjectCompareDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_object_compare_detail.test"
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
				Config: testDataSourceObjectCompareDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "compare_details.#"),
				),
			},
		},
	})
}

// This resource has no response data, and the filtering parameters cannot be validated in the test case.
func testDataSourceObjectCompareDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_object_compare_detail" "test" {
  job_id       = "%s"
  compare_type = "TABLE"
}
`, acceptance.HW_DRS_JOB_IDS)
}
