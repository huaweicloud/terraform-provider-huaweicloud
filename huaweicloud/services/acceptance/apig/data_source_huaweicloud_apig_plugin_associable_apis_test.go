package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataPluginAssociableApis_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_plugin_associable_apis.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byApiName   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_api_name"
		dcByApiName = acceptance.InitDataSourceCheck(byApiName)

		byApiId   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_api_id"
		dcByApiId = acceptance.InitDataSourceCheck(byApiId)

		byGroupId   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_group_id"
		dcByGroupId = acceptance.InitDataSourceCheck(byGroupId)

		byReqMethod   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_req_method"
		dcByReqMethod = acceptance.InitDataSourceCheck(byReqMethod)

		byTags   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)

		byNoTags   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_no_tags"
		dcByNoTags = acceptance.InitDataSourceCheck(byNoTags)

		byNotFoundApiName   = "data.huaweicloud_apig_plugin_associable_apis.filter_by_not_found_api_name"
		dcByNotFoundApiName = acceptance.InitDataSourceCheck(byNotFoundApiName)

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
				Config: testAccDataPluginAssociableApis_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "apis.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(rName, "apis.0.id"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.name"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.type"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.req_protocol"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.req_method"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.req_uri"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.auth_type"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.match_mode"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.group_id"),
					resource.TestCheckResourceAttrSet(rName, "apis.0.group_name"),
					dcByApiName.CheckResourceExists(),
					resource.TestCheckOutput("is_api_name_filter_useful", "true"),
					dcByNotFoundApiName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_api_name_useful", "true"),
					dcByApiId.CheckResourceExists(),
					resource.TestCheckOutput("is_api_id_filter_useful", "true"),
					dcByGroupId.CheckResourceExists(),
					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
					dcByReqMethod.CheckResourceExists(),
					resource.TestCheckOutput("is_req_method_filter_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByNoTags.CheckResourceExists(),
					resource.TestCheckOutput("is_no_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataPluginAssociableApis_base(name string) string {
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
  count = 2

  instance_id             = local.instance_id
  group_id                = huaweicloud_apig_group.test.id
  name                    = format("%[3]s_%%s", count.index)
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
    name             = format("%[3]s_web_%%s", count.index)
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

  tags = count.index == 0 ? ["terraform", "test"] : null

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
  count = 2

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test[count.index].id
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

func testAccDataPluginAssociableApis_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_plugin_associable_apis" "test" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

# Filter by API name (fuzzy search)
locals {
  api_name = "%[2]s"
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_api_name" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_name    = local.api_name

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  api_name_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_api_name.apis : strcontains(v.name, local.api_name)
  ]
}

output "is_api_name_filter_useful" {
  value = length(local.api_name_filter_result) >= 2 && alltrue(local.api_name_filter_result)
}

# Filter by not found API name
locals {
  not_found_api_name = "not_found_api"
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_not_found_api_name" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_name    = local.not_found_api_name

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  not_found_api_name_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_not_found_api_name.apis : strcontains(v.name, local.not_found_api_name)
  ]
}

output "is_not_found_api_name_useful" {
  value = length(local.not_found_api_name_filter_result) == 0
}

# Filter by API ID
locals {
  api_id = huaweicloud_apig_api.test[0].id
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_api_id" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  api_id      = local.api_id

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  api_id_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_api_id.apis : v.id == local.api_id
  ]
}

output "is_api_id_filter_useful" {
  value = length(local.api_id_filter_result) > 0 && alltrue(local.api_id_filter_result)
}

# Filter by group ID
locals {
  group_id = huaweicloud_apig_group.test.id
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_group_id" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  group_id    = local.group_id

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  group_id_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_group_id.apis : v.group_id == local.group_id
  ]
}

output "is_group_id_filter_useful" {
  value = length(local.group_id_filter_result) > 0 && alltrue(local.group_id_filter_result)
}

# Filter by request method
locals {
  req_method = huaweicloud_apig_api.test[0].request_method
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_req_method" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  req_method  = local.req_method

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  req_method_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_req_method.apis : v.req_method == local.req_method
  ]
}

output "is_req_method_filter_useful" {
  value = length(local.req_method_filter_result) > 0 && alltrue(local.req_method_filter_result)
}

# Filter by tags
locals {
  tags = "terraform"
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_tags" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  tags        = local.tags

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_tags.apis : contains(v.tags, local.tags)
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}

# Filter by tags (no tags)
locals {
  no_tags = "#no_tags#"
}

data "huaweicloud_apig_plugin_associable_apis" "filter_by_no_tags" {
  instance_id = local.instance_id
  plugin_id   = huaweicloud_apig_plugin.test.id
  env_id      = huaweicloud_apig_environment.test.id
  tags        = local.no_tags

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}

locals {
  no_tags_filter_result = [
    for v in data.huaweicloud_apig_plugin_associable_apis.filter_by_no_tags.apis : length(v.tags) == 0
  ]
}

output "is_no_tags_filter_useful" {
  value = length(local.no_tags_filter_result) > 0 && alltrue(local.no_tags_filter_result)
}
`, testAccDataPluginAssociableApis_base(name), name)
}
