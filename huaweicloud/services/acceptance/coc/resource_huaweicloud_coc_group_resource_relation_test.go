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
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
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
  depends_on = [huaweicloud_coc_group.test]

  cloud_service_name = "ecs"
  type               = "cloudservers"
  resource_id_list   = [huaweicloud_compute_instance.test.id]
}

resource "time_sleep" "wait_10_seconds" {
  depends_on = [huaweicloud_coc_group.test]

  create_duration = "10s"
}

resource "huaweicloud_coc_group_resource_relation" "test" {
  depends_on = [time_sleep.wait_10_seconds]

  group_id         = huaweicloud_coc_group.test.id
  cmdb_resource_id = data.huaweicloud_coc_resources.test.data[0].id
}
`, testAccComputeInstance_basic(name), testAccGroup_basic(name))
}
