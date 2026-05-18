package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccJobConfigurationUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testJobConfigurationUpdate_basic(),
			},
		},
	})
}

func testJobConfigurationUpdate_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_job_configurations" "test" {
  job_id = "%[1]s"
}

resource "huaweicloud_drs_job_configuration_update" "test" {
  job_id = "%[1]s"

  values {
    parameter_name  = data.huaweicloud_drs_job_configurations.test.parameter_config_list.0.name
    parameter_value = data.huaweicloud_drs_job_configurations.test.parameter_config_list.0.value
  }
}
`, acceptance.HW_DRS_JOB_ID)
}
