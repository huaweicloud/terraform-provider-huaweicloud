package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccApiVersionUnpublish_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		rcName = "huaweicloud_apig_api_version_unpublish.test"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApiVersionUnpublish_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rcName, "id"),
				),
			},
		},
	})
}

func testAccApiVersionUnpublish_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

resource "huaweicloud_apig_group" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name        = "%[3]s"
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name        = "%[3]s"
}

resource "huaweicloud_apig_api" "test" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  group_id         = huaweicloud_apig_group.test.id
  name             = "%[3]s"
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = "/mock/test"

  mock {
    status_code = 200
  }
}

resource "huaweicloud_apig_api_action" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  action      = "online"
}

resource "huaweicloud_apig_api_version_unpublish" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  version_id  = huaweicloud_apig_api_action.test.version_id
}
`, common.TestBaseNetwork(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
