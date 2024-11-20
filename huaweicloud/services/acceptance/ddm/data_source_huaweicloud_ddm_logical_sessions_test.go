package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdmLogicalSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_logical_sessions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDMInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdmLogicalSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.user"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.db"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.command"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_processes.0.info"),

					resource.TestCheckOutput("keyword_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDdmLogicalSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ddm_logical_sessions" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_ddm_logical_sessions" "keyword_filter" {
  instance_id = "%[1]s"
  keyword     = "SELECT"
}
output "keyword_filter_is_useful" {
  value = length(data.huaweicloud_ddm_logical_sessions.keyword_filter.logical_processes) > 0
}
`, acceptance.HW_DDM_INSTANCE_ID)
}
