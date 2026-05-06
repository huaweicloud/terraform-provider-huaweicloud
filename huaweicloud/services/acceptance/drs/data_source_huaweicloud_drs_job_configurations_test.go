package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsJobConfigurations_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_job_configurations.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsJobConfigurations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameter_config_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameter_config_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parameter_config_list.0.value"),
				),
			},
		},
	})
}

func testAccDataSourceDrsJobConfigurations_basic() string {
	return fmt.Sprintf(`data "huaweicloud_drs_job_configurations" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}
