package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterPorts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cpcs_cluster_ports.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCpcsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterPorts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "result.#"),
				),
			},
		},
	})
}

func testAccDataSourceClusterPorts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cpcs_cluster_ports" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_CPCS_CLUSTER_ID)
}
