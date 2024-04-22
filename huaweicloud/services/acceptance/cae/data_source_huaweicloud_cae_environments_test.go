package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceEnvironments_basic(t *testing.T) {
	rName := "data.huaweicloud_cae_environments.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceEnvironments_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckOutput("environment_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceEnvironments_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cae_environments" "test" {
  environment_id = "%[1]s"
}

locals {
  environment_id_filter_result = [for v in data.huaweicloud_cae_environments.test.environments : v.id == "%[1]s"]
}

output "environment_id_filter_is_useful" {
  value = alltrue(local.environment_id_filter_result) && length(local.environment_id_filter_result) > 0
}

data "huaweicloud_cae_environments" "test1" {
  status = "finish"
}

locals {
  status_filter_result = [for v in data.huaweicloud_cae_environments.test1.environments : v.status == "finish"]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, acceptance.HW_CAE_ENVIRONMENT_ID)
}
