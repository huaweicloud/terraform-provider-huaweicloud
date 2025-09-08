package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceApiHistoryVersions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_apig_api_history_versions.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		dcFilterByEnvIdName = "data.huaweicloud_apig_api_history_versions.filter_by_env_id"
		dcFilterByEnvId     = acceptance.InitDataSourceCheck(dcFilterByEnvIdName)

		dcFilterByEnvNameName = "data.huaweicloud_apig_api_history_versions.filter_by_env_name"
		dcFilterByEnvName     = acceptance.InitDataSourceCheck(dcFilterByEnvNameName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiHistoryVersions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "api_versions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.number"),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.api_id"),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.env_id"),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.env_name"),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.publish_time"),
					resource.TestCheckResourceAttrSet(dcName, "api_versions.0.status"),
					dcFilterByEnvId.CheckResourceExists(),
					resource.TestCheckOutput("is_env_id_filter_useful", "true"),
					dcFilterByEnvName.CheckResourceExists(),
					resource.TestCheckOutput("is_env_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceApiHistoryVersions_base(name string) string {
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
  count = 2

  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name        = format("%[3]s_%%d", count.index)
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

resource "huaweicloud_apig_api_action" "test_online_first" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[0].id
  action      = "online"
}

resource "huaweicloud_apig_api_action" "test_offline_first" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[0].id
  action      = "offline"
  
  depends_on = [huaweicloud_apig_api_action.test_online_first]
}

resource "huaweicloud_apig_api_action" "test_online_second" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[1].id
  action      = "online"
  
  depends_on = [huaweicloud_apig_api_action.test_offline_first]
}

resource "huaweicloud_apig_api_action" "test_offline_second" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[1].id
  action      = "offline"
  
  depends_on = [huaweicloud_apig_api_action.test_online_second]
}
`, common.TestBaseNetwork(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceApiHistoryVersions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all
data "huaweicloud_apig_api_history_versions" "test" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id

  depends_on = [
    huaweicloud_apig_api_action.test_online_first,
    huaweicloud_apig_api_action.test_online_second,
  ]
}

# Filter by api environment id
data "huaweicloud_apig_api_history_versions" "filter_by_env_id" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[0].id

  depends_on = [
    huaweicloud_apig_api_action.test_online_first,
    huaweicloud_apig_api_action.test_online_second,
  ]
}

output "is_env_id_filter_useful" {
  value = length(data.huaweicloud_apig_api_history_versions.filter_by_env_id.api_versions) > 0 && alltrue(
    [for v in data.huaweicloud_apig_api_history_versions.filter_by_env_id.api_versions[*].env_id :
v == huaweicloud_apig_environment.test[0].id]
  )
}

# Filter by api environment name
data "huaweicloud_apig_api_history_versions" "filter_by_env_name" {
  instance_id = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  api_id      = huaweicloud_apig_api.test.id
  env_name    = huaweicloud_apig_environment.test[1].name

  depends_on = [
    huaweicloud_apig_api_action.test_online_first,
    huaweicloud_apig_api_action.test_online_second,
	data.huaweicloud_apig_api_history_versions.filter_by_env_id
  ]
}

output "is_env_name_filter_useful" {
  value = length(data.huaweicloud_apig_api_history_versions.filter_by_env_name.api_versions) > 0 && alltrue(
    [for v in data.huaweicloud_apig_api_history_versions.filter_by_env_name.api_versions[*].env_name :
v == huaweicloud_apig_environment.test[1].name]
  )
}
`, testAccDataSourceApiHistoryVersions_base(name))
}
