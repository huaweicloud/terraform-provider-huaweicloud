package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceGroupResourceRelation_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupResourceRelation_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAccGroupResourceRelation_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_coc_resources" "test" {
  cloud_service_name = "ecs"
  type               = "cloudservers"
  resource_id_list   = [huaweicloud_compute_instance.test.id]

  depends_on         = [huaweicloud_coc_group.test]
}

locals {
  cmdb_resource_id_list = data.huaweicloud_coc_resources.test.data[*].id
}

resource "huaweicloud_coc_group_resource_relation" "test" {
  group_id              = huaweicloud_coc_group.test.id
  cmdb_resource_id_list = local.cmdb_resource_id_list
}
`, testAccComputeInstance_basic(name), testAccGroup_basic(name))
}
