package cse

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMicroserviceEngines_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_cse_microservice_engines.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSEMicroserviceEngineID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMicroserviceEngines_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("filter_engine_by_id", "true"),
				),
			},
		},
	})
}

func testAccDataMicroserviceEngines_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cse_microservice_engines" "test" {}

locals {
  id_filter_result = [
	for o in data.huaweicloud_cse_microservice_engines.test.engines : o if o.id == "%[1]s"
  ]
}

output "filter_engine_by_id" {
  value = length(local.id_filter_result) > 0
}
`, acceptance.HW_CSE_MICROSERVICE_ENGINE_ID)
}
