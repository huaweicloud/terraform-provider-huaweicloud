package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocGroupSync_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocGroupSync_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocGroupSync_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_coc_group_sync" "test" {
  group_id           = huaweicloud_coc_group.test.id
  cloud_service_name = "ecs"
  type               = "cloudservers"
}
`, testAccGroup_basic(name))
}
