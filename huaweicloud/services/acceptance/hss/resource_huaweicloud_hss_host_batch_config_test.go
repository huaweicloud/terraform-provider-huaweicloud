package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHostBatchConfig_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid HSS host protection host id and config it to the environment variable.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testHostBatchConfig_basic(),
			},
		},
	})
}

func testHostBatchConfig_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_host_batch_config" "test" {
  operate_all = false
  host_ids    = ["%s"]
  mode        = "default"
  cpu_limit   = "0.2"
  mem_limit   = "800"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
