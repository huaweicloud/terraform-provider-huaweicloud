package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataApplicationAuthorizeStatistic_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_apig_application_authorize_statistic.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationAuthorizeStatistic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "authed_nums", regexp.MustCompile(`^[1-9]+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "unauthed_nums", regexp.MustCompile(`^[1-9]+$`)),
				),
			},
		},
	})
}

func testAccDataApplicationAuthorizeStatistic_basic() string {
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
  api_id      = huaweicloud_apig_api.test[0].id
  env_id      = huaweicloud_apig_environment.test.id
}

resource "huaweicloud_apig_application" "test" {
  instance_id = local.instance_id
  name        = "%[2]s"
}

resource "huaweicloud_apig_application_authorization" "test" {
  instance_id    = local.instance_id
  application_id = huaweicloud_apig_application.test.id
  env_id         = huaweicloud_apig_environment.test.id
  api_ids        = slice(huaweicloud_apig_api.test[*].id, 0, 1)

  depends_on = [
    huaweicloud_apig_api_publishment.test,
  ]
}


data "huaweicloud_apig_application_authorize_statistic" "test" {
  instance_id = local.instance_id

  depends_on = [
    huaweicloud_apig_application_authorization.test,
  ]
}
`, common.TestBaseComputeResources(name),
		acceptance.HW_APIG_DEDICATED_INSTANCE_ID,
		name,
		acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}
