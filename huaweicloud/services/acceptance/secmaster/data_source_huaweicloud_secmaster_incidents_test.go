package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIncidents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_incidents.test"
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
				Config: testDataSourceIncidents_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "incidents.#"),
					resource.TestCheckResourceAttrSet(dataSource, "incidents.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "incidents.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "incidents.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "incidents.0.status"),

					resource.TestCheckOutput("condition_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceIncidents_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_incidents" "test" {
  workspace_id = "%[2]s"

  condition {
    conditions {
      name = "severity"
      data = [ "severity", "=", "Tips" ]
    }
    logics = ["severity"]
  }
  
  depends_on = [huaweicloud_secmaster_incident.test]
}

output "condition_filter_is_useful" {
  value = length(data.huaweicloud_secmaster_incidents.test.incidents) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_incidents.test.incidents[*].level : v == "Tips"]
  )
}
`, testIncident_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID)
}
