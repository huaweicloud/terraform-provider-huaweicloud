package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceApplications_basic(t *testing.T) {
	rName := "data.huaweicloud_cae_applications.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceApplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckOutput("application_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceApplications_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cae_applications" "test" {
  environment_id = "%[1]s"
  application_id = "%[2]s"
}

locals {
  application_id_filter_result = [for v in data.huaweicloud_cae_applications.test.applications : v.id == "%[2]s"]
}

output "application_id_filter_is_useful" {
  value = alltrue(local.application_id_filter_result) && length(local.application_id_filter_result) > 0
}
`, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID)
}
