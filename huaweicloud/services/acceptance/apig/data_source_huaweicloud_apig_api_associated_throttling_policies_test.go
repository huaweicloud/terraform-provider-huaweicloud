package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceApiAssociatedThrottlingPolicies_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_api_associated_throttling_policies.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_apig_api_associated_throttling_policies.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_api_associated_throttling_policies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_api_associated_throttling_policies.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byType   = "data.huaweicloud_apig_api_associated_throttling_policies.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byEnvName   = "data.huaweicloud_apig_api_associated_throttling_policies.filter_by_env_name"
		dcByEnvName = acceptance.InitDataSourceCheck(byEnvName)

		byPeriodUnit   = "data.huaweicloud_apig_api_associated_throttling_policies.filter_by_period_unit"
		dcByPeriodUnit = acceptance.InitDataSourceCheck(byPeriodUnit)
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
				Config: testAccDataSourceApiAssociatedThrottlingPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(rName, "policies.0.max_api_requests", "100"),
					resource.TestMatchResourceAttr(rName, "policies.0.user_throttles.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(rName, "policies.0.app_throttles.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrPair(rName, "policies.0.user_throttles.0.throttling_object_id",
						"huaweicloud_apig_throttling_policy.test", "user_throttles.0.throttling_object_id"),
					resource.TestCheckResourceAttrPair(rName, "policies.0.app_throttles.0.throttling_object_id",
						"huaweicloud_apig_throttling_policy.test", "app_throttles.0.throttling_object_id"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByEnvName.CheckResourceExists(),
					resource.TestCheckOutput("is_env_name_filter_useful", "true"),
					dcByPeriodUnit.CheckResourceExists(),
					resource.TestCheckOutput("is_period_unit_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceApiAssociatedThrottlingPolicies_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[3]s"
  }
}

data "huaweicloud_identity_users" "test" {}

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[4]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = local.instance_id
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  port             = 8000
  balance_strategy = 2
  member_type        = "ecs"
  type               = 2

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
  name                    = "%[2]s"
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
    name             = "%[2]s_policy1"
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
}

resource "huaweicloud_apig_environment" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "huaweicloud_apig_api_publishment" "test" {
  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test.id
}

resource "huaweicloud_apig_application" "test" {
  name        = "%[2]s"
  instance_id = local.instance_id
}

resource "huaweicloud_apig_throttling_policy" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  type             = "API-based"
  period           = 15000
  period_unit      = "SECOND"
  max_api_requests = 100

  app_throttles {
    max_api_requests     = 30
    throttling_object_id = huaweicloud_apig_application.test.id
  }
  user_throttles {
    max_api_requests     = 30
    throttling_object_id = data.huaweicloud_identity_users.test.users[0].id
  }
}

resource "huaweicloud_apig_throttling_policy_associate" "test" {
  instance_id = local.instance_id
  policy_id   = huaweicloud_apig_throttling_policy.test.id
  publish_ids = [
    huaweicloud_apig_api_publishment.test.publish_id
  ]
}
`, common.TestBaseComputeResources(name), name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID,
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccDataSourceApiAssociatedThrottlingPolicies_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_api_associated_throttling_policies" "test" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
}

# Filter by ID
locals {
  policy_id = huaweicloud_apig_throttling_policy.test.id
}

data "huaweicloud_apig_api_associated_throttling_policies" "filter_by_id" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  policy_id = local.policy_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_api_associated_throttling_policies.filter_by_id.policies[*].id : v == local.policy_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  policy_name = huaweicloud_apig_throttling_policy.test.name
}

data "huaweicloud_apig_api_associated_throttling_policies" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  name = local.policy_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_api_associated_throttling_policies.filter_by_name.policies[*].name : v == local.policy_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_apig_api_associated_throttling_policies" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  name = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.huaweicloud_apig_api_associated_throttling_policies.filter_by_not_found_name.policies[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}

# Filter by type
locals {
  policy_type = huaweicloud_apig_throttling_policy.test.type
}

data "huaweicloud_apig_api_associated_throttling_policies" "filter_by_type" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  type = local.policy_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_api_associated_throttling_policies.filter_by_type.policies[*].type : v == local.policy_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by env name
locals {
  env_name = huaweicloud_apig_environment.test.name
}

data "huaweicloud_apig_api_associated_throttling_policies" "filter_by_env_name" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  env_name = local.env_name
}

locals {
  env_name_filter_result = [
    for v in data.huaweicloud_apig_api_associated_throttling_policies.filter_by_env_name.policies[*].env_name : v == local.env_name
  ]
}

output "is_env_name_filter_useful" {
  value = length(local.env_name_filter_result) > 0 && alltrue(local.env_name_filter_result)
}

# Filter by period unit
locals {
  period_unit = huaweicloud_apig_throttling_policy.test.period_unit
}

data "huaweicloud_apig_api_associated_throttling_policies" "filter_by_period_unit" {
  depends_on = [
    huaweicloud_apig_throttling_policy_associate.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  period_unit = local.period_unit
}

locals {
  period_unit_filter_result = [
    for v in data.huaweicloud_apig_api_associated_throttling_policies.filter_by_period_unit.policies[*].period_unit : v == local.period_unit
  ]
}

output "is_period_unit_filter_useful" {
  value = length(local.period_unit_filter_result) > 0 && alltrue(local.period_unit_filter_result)
}
`, testAccDataSourceApiAssociatedThrottlingPolicies_base())
}
