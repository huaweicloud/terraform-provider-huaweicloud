package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccClusterProtectSwitchMode_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test requires setting a cluster ID under the default enterprise project that has been connected to
			// the HSS service.
			acceptance.TestAccPreCheckHSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testClusterProtectSwitchMode_basic(),
			},
		},
	})
}

func testClusterProtectSwitchMode_basic() string {
	return fmt.Sprintf(`

resource "huaweicloud_hss_cluster_protect_switch_mode" "test" {
  cluster_ids           = ["%[1]s"]
  opr                   = 1
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_CLUSTER_ID)
}
