package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAgentChecks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_cbr_agent_checks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Running the test, prepare an ECS instance, and enable backup.
			acceptance.TestAccPreCheckCbrResourceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAgentChecks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_status_attr.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_status_attr.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_status_attr.0.installed"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_status_attr.0.code"),
				),
			},
		},
	})
}

func testAccDataSourceAgentChecks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbr_agent_checks" "test" {
  agent_status {
    resource_id   = "%s"
    resource_type = "OS::Nova::Server"
  }
}
`, acceptance.HW_CBR_RESOURCE_ID)
}
