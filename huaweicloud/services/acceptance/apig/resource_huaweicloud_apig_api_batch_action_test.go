package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccApigApiBatchAction_basic(t *testing.T) {
	var (
		name                  = acceptance.RandomAccResourceName()
		rcWithOnlineName      = "huaweicloud_apig_api_batch_action.batch_online_apis_for_env"
		rcWithGroupOnlineName = "huaweicloud_apig_api_batch_action.batch_online_apis_for_group"
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
				Config: testApigApiBatchAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rcWithOnlineName, "action", "online"),
					resource.TestCheckResourceAttr(rcWithOnlineName, "apis.#", "2"),
					resource.TestCheckResourceAttrSet(rcWithOnlineName, "id"),
				),
			},
			{
				Config: testApigApiBatchAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rcWithGroupOnlineName, "action", "online"),
					resource.TestCheckResourceAttrSet(rcWithGroupOnlineName, "group_id"),
					resource.TestCheckResourceAttrSet(rcWithGroupOnlineName, "id"),
				),
			},
		},
	})
}

func testApigApiBatchAction_base(name string) string {
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
  count = 2

  instance_id      = huaweicloud_apig_instance.test.id
  group_id         = huaweicloud_apig_group.test.id
  name             = format("%[2]s_%%d", count.index)
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = format("/mock/test%%d", count.index)

  mock {
    status_code = 200
  }
}
`, common.TestBaseNetwork(name), name)
}

func testApigApiBatchAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api_batch_action" "batch_online_apis_for_env" {
  instance_id = huaweicloud_apig_instance.test.id
  action      = "online"
  env_id      = huaweicloud_apig_environment.test.id
  remark      = "Test batch action"
  apis        = huaweicloud_apig_api.test[*].id
}

resource "huaweicloud_apig_api_batch_action" "batch_offline_apis_for_env" {
  instance_id = huaweicloud_apig_instance.test.id
  action      = "offline"
  env_id      = huaweicloud_apig_environment.test.id
  remark      = "Test batch action"
  apis        = huaweicloud_apig_api.test[*].id

  depends_on = [
    huaweicloud_apig_api_batch_action.batch_online_apis_for_env,
  ]
}
`, testApigApiBatchAction_base(name))
}

func testApigApiBatchAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api_batch_action" "batch_online_apis_for_group" {
  instance_id = huaweicloud_apig_instance.test.id
  action      = "online"
  env_id      = huaweicloud_apig_environment.test.id
  group_id    = huaweicloud_apig_group.test.id
  remark      = "Test batch action by group"

  depends_on = [
    huaweicloud_apig_api.test,
  ]
}

resource "huaweicloud_apig_api_batch_action" "batch_offline_apis_for_group" {
  instance_id = huaweicloud_apig_instance.test.id
  action      = "offline"
  env_id      = huaweicloud_apig_environment.test.id
  group_id    = huaweicloud_apig_group.test.id
  remark      = "Test batch action by group"

  depends_on = [
    huaweicloud_apig_api.test,
    huaweicloud_apig_api_batch_action.batch_online_apis_for_group,
  ]
}
`, testApigApiBatchAction_base(name))
}
