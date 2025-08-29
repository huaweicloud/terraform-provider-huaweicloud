package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocOtherResourceUniAgentSync_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocOtherResourceUniAgentSync_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocOtherResourceUniAgentSync_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_coc_other_resource_uniagent_sync" "test" {
  resource_infos {
    region_id   = huaweicloud_compute_instance.test.region
    resource_id = huaweicloud_compute_instance.test.id
  }
}
`, testAccComputeInstance_basic(name))
}
