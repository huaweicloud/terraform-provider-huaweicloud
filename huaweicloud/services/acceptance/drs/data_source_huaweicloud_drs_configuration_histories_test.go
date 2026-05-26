package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsConfigurationHistories_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_configuration_histories.test"
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
				Config: testDataSourceDrsConfigurationHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameter_history_config_list.#"),
				),
			},
		},
	})
}

func testDataSourceDrsConfigurationHistories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_configuration_histories" "test" {
  job_id     = "%s"
  begin_time = "2020-09-01T18:50:20Z"
  end_time   = "2025-09-01T18:50:20Z"
}
`, acceptance.HW_DRS_JOB_ID)
}
