package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataApiAssociablePlugins_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_api_associable_plugins.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byEnvId   = "data.huaweicloud_apig_api_associable_plugins.filter_by_env_id"
		dcByEnvId = acceptance.InitDataSourceCheck(byEnvId)

		byPluginName   = "data.huaweicloud_apig_api_associable_plugins.filter_by_plugin_name"
		dcByPluginName = acceptance.InitDataSourceCheck(byPluginName)

		byNotFoundPluginName   = "data.huaweicloud_apig_api_associable_plugins.filter_by_not_found_plugin_name"
		dcByNotFoundPluginName = acceptance.InitDataSourceCheck(byNotFoundPluginName)

		byPluginType   = "data.huaweicloud_apig_api_associable_plugins.filter_by_plugin_type"
		dcByPluginType = acceptance.InitDataSourceCheck(byPluginType)

		byPluginId   = "data.huaweicloud_apig_api_associable_plugins.filter_by_plugin_id"
		dcByPluginId = acceptance.InitDataSourceCheck(byPluginId)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApiAssociablePlugins_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "plugins.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.id"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.name"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.type"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.scope"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.content"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.description"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.create_time"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.update_time"),
					dcByEnvId.CheckResourceExists(),
					resource.TestCheckOutput("is_env_id_filter_useful", "true"),
					dcByPluginName.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_name_filter_useful", "true"),
					dcByNotFoundPluginName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_plugin_name_useful", "true"),
					dcByPluginType.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_type_filter_useful", "true"),
					dcByPluginId.CheckResourceExists(),
					resource.TestCheckOutput("is_plugin_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataApiAssociablePlugins_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[3]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[4]s"
  }
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[3]s"
  instance_id = local.instance_id
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[3]s"
  port             = 8000
  balance_strategy = 2
  member_type      = "ecs"
  type             = 2

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 10  # maximum value
    threshold_abnormal = 10  # maximum value
    interval           = 300 # maximum value
    timeout            = 30  # maximum value
    path               = "/"
    method             = "HEAD"
    port               = 8080
    http_codes         = "201,202,303-404"
  }

  member {
    id   = huaweicloud_compute_instance.test.id
    name = huaweicloud_compute_instance.test.name
  }
}

resource "huaweicloud_apig_api" "test" {
  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[3]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/{user_age}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }

  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = huaweicloud_apig_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }

  web_policy {
    name             = "%[3]s_web"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getUserAge/{userAge}"
    timeout          = 30000
    vpc_channel_id   = huaweicloud_apig_channel.test.id

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }

    conditions {
      source     = "param"
      param_name = "user_age"
      type       = "Equal"
      value      = "28"
    }
  }

  lifecycle {
    ignore_changes = [
      request_params,
    ]
  }
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = local.instance_id
  name        = "%[3]s"
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
}

resource "huaweicloud_apig_plugin" "test" {
  instance_id = local.instance_id
  name        = "%[3]s_cors"
  type        = "cors"
  description = "Created by terraform script"
  content     = jsonencode({
    allow_origin      = "*"
    allow_methods     = "GET,PUT,DELETE,HEAD,PATCH"
    allow_headers     = "Content-Type,Accept,Cache-Control"
    expose_headers    = "X-Request-Id,X-Apig-Latency"
    max_age           = 12700
    allow_credentials = true
  })
}
`, common.TestBaseComputeResources(name),
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccDataApiAssociablePlugins_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_api_associable_plugins" "test" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

# Filter by environment ID
locals {
  env_id = huaweicloud_apig_environment.test.id
}

data "huaweicloud_apig_api_associable_plugins" "filter_by_env_id" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = local.env_id

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  env_id_filter_result = [
    for v in data.huaweicloud_apig_api_associable_plugins.filter_by_env_id.plugins : v != null
  ]
}

output "is_env_id_filter_useful" {
  value = length(local.env_id_filter_result) > 0
}

# Filter by plugin name (fuzzy search)
locals {
  plugin_name = "%[2]s"
}

data "huaweicloud_apig_api_associable_plugins" "filter_by_plugin_name" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_name = local.plugin_name

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  plugin_name_filter_result = [
    for v in data.huaweicloud_apig_api_associable_plugins.filter_by_plugin_name.plugins : strcontains(v.name, local.plugin_name)
  ]
}

output "is_plugin_name_filter_useful" {
  value = length(local.plugin_name_filter_result) > 0 && alltrue(local.plugin_name_filter_result)
}

# Filter by not found plugin name
locals {
  not_found_plugin_name = "not_found_plugin"
}

data "huaweicloud_apig_api_associable_plugins" "filter_by_not_found_plugin_name" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_name = local.not_found_plugin_name

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  not_found_plugin_name_filter_result = [
    for v in data.huaweicloud_apig_api_associable_plugins.filter_by_not_found_plugin_name.plugins : strcontains(v.name, local.not_found_plugin_name)
  ]
}

output "is_not_found_plugin_name_useful" {
  value = length(local.not_found_plugin_name_filter_result) == 0
}

# Filter by plugin type
locals {
  plugin_type = "cors"
}

data "huaweicloud_apig_api_associable_plugins" "filter_by_plugin_type" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_type = local.plugin_type

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  plugin_type_filter_result = [
    for v in data.huaweicloud_apig_api_associable_plugins.filter_by_plugin_type.plugins : v.type == local.plugin_type
  ]
}

output "is_plugin_type_filter_useful" {
  value = length(local.plugin_type_filter_result) > 0 && alltrue(local.plugin_type_filter_result)
}

# Filter by plugin ID
locals {
  plugin_id = huaweicloud_apig_plugin.test.id
}

data "huaweicloud_apig_api_associable_plugins" "filter_by_plugin_id" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
  plugin_id   = local.plugin_id

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  plugin_id_filter_result = [
    for v in data.huaweicloud_apig_api_associable_plugins.filter_by_plugin_id.plugins : v.id == local.plugin_id
  ]
}

output "is_plugin_id_filter_useful" {
  value = length(local.plugin_id_filter_result) > 0 && alltrue(local.plugin_id_filter_result)
}
`, testAccDataApiAssociablePlugins_base(name), name)
}
