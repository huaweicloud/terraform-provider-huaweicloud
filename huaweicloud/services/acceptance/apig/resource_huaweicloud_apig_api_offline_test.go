package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccApiOffline_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApiOffline_basic(name),
			},
		},
	})
}

func testAccApiOffline_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_apig_instance" "test" {
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), [])
  edition               = "BASIC"
  name                  = "%[2]s"
  description           = "created by acc test for API offline action"
}

resource "huaweicloud_apig_group" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_apig_api" "test" {
  instance_id      = huaweicloud_apig_instance.test.id
  group_id         = huaweicloud_apig_group.test.id
  name             = "%[2]s"
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = "/mock/test"

  mock {
    status_code = 200
  }
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccApiOffline_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api_offline" "test" {
  instance_id = huaweicloud_apig_instance.test.id
  version_id  = try(huaweicloud_apig_api_publishment.test.histories[0].version_id, "")

  depends_on = [huaweicloud_apig_api_publishment.test]
}

# Recover api status for api rollback action
resource "huaweicloud_apig_api_publishment" "recover" {
  instance_id = huaweicloud_apig_instance.test.id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id

  depends_on = [huaweicloud_apig_api_offline.test]
}
`, testAccApiOffline_base(name))
}
