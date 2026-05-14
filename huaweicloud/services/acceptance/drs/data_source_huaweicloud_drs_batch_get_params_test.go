package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsBatchGetParams_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_batch_get_params.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsBatchGetParams_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.0.params.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.0.params.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.0.params.0.group"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.0.params.0.source_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "params_list.0.params.0.target_value"),
				),
			},
		},
	})
}

func testAccDataSourceDrsBatchGetParams_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_batch_get_params" "test" {
  job_ids = split(",", "%s")
  refresh = 1
}
`, acceptance.HW_DRS_JOB_IDS)
}
