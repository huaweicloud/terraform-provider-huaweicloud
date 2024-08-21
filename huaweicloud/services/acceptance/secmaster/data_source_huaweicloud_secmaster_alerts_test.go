package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlerts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_alerts.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlerts_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alerts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "alerts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "alerts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alerts.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "alerts.0.status"),

					resource.TestCheckOutput("condition_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAlerts_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_alerts" "test" {
  workspace_id = "%[2]s"

  condition {
    conditions {
      name = "severity"
      data = [ "severity", "=", "Tips" ]
    }
    logics = ["severity"]
  }
  
  depends_on = [huaweicloud_secmaster_alert.test]
}

output "condition_filter_is_useful" {
  value = length(data.huaweicloud_secmaster_alerts.test.alerts) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_alerts.test.alerts[*].level : v == "Tips"]
  )
}
`, testAlert_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}
