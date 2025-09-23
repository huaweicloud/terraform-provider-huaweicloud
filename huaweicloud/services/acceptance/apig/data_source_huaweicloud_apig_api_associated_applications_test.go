package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceApiAssociatedApplications_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_apig_api_associated_applications.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_apig_api_associated_applications.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_api_associated_applications.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_api_associated_applications.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byEnvId   = "data.huaweicloud_apig_api_associated_applications.filter_by_env_id"
		dcByEnvId = acceptance.InitDataSourceCheck(byEnvId)

		byEnvName   = "data.huaweicloud_apig_api_associated_applications.filter_by_env_name"
		dcByEnvName = acceptance.InitDataSourceCheck(byEnvName)
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
				Config: testAccDataSourceApiAssociatedApplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(rName, "applications.0.id"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.name"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.description"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.env_id"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.env_name"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.bind_id"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.bind_time"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByEnvId.CheckResourceExists(),
					resource.TestCheckOutput("is_env_id_filter_useful", "true"),
					dcByEnvName.CheckResourceExists(),
					resource.TestCheckOutput("is_env_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceApiAssociatedApplications_base() string {
	name := acceptance.RandomAccResourceName()

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
}

resource "huaweicloud_apig_environment" "test" {
  count = 2

  instance_id = local.instance_id
  name        = "%[3]s_${count.index}"
}

resource "huaweicloud_apig_api_publishment" "test" {
  count = 2

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
  env_id      = huaweicloud_apig_environment.test[count.index].id
}

resource "huaweicloud_apig_application" "test" {
  instance_id = local.instance_id
  name        = "%[3]s"
  description = "Created by terraform script"
}

resource "huaweicloud_apig_application_authorization" "test" {
  depends_on = [huaweicloud_apig_api_publishment.test]

  count = 2

  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  env_id         = huaweicloud_apig_environment.test[count.index].id
  api_ids        = [huaweicloud_apig_api.test.id]
}
`, common.TestBaseComputeResources(name),
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccDataSourceApiAssociatedApplications_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_api_associated_applications" "test" {
  depends_on = [
    huaweicloud_apig_application_authorization.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id
}

# Filter by ID
locals {
  application_id = huaweicloud_apig_application.test.id
}

data "huaweicloud_apig_api_associated_applications" "filter_by_id" {
  depends_on = [
    huaweicloud_apig_application_authorization.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  application_id = local.application_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_api_associated_applications.filter_by_id.applications[*].id : v == local.application_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  application_name = huaweicloud_apig_application.test.name
}

data "huaweicloud_apig_api_associated_applications" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_application_authorization.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  name = local.application_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_api_associated_applications.filter_by_name.applications[*].name : v == local.application_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_apig_api_associated_applications" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_apig_application_authorization.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  name = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.huaweicloud_apig_api_associated_applications.filter_by_not_found_name.applications[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}

# Filter by env ID
locals {
  env_id = huaweicloud_apig_application_authorization.test[0].env_id
}

data "huaweicloud_apig_api_associated_applications" "filter_by_env_id" {
  depends_on = [
    huaweicloud_apig_application_authorization.test,
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  env_id = local.env_id
}

locals {
  env_id_filter_result = [
    for v in data.huaweicloud_apig_api_associated_applications.filter_by_env_id.applications[*].env_id : v == local.env_id
  ]
}

output "is_env_id_filter_useful" {
  value = length(local.env_id_filter_result) > 0 && alltrue(local.env_id_filter_result)
}

# Filter by env name
locals {
  env_name = huaweicloud_apig_environment.test[0].name
}

data "huaweicloud_apig_api_associated_applications" "filter_by_env_name" {
  depends_on = [
    huaweicloud_apig_application_authorization.test
  ]

  instance_id = local.instance_id
  api_id      = huaweicloud_apig_api.test.id

  env_name = local.env_name
}

locals {
  env_name_filter_result = [
    for v in data.huaweicloud_apig_api_associated_applications.filter_by_env_name.applications[*].env_name : v == local.env_name
  ]
}

output "is_env_name_filter_useful" {
  value = length(local.env_name_filter_result) > 0 && alltrue(local.env_name_filter_result)
}
`, testAccDataSourceApiAssociatedApplications_base())
}
