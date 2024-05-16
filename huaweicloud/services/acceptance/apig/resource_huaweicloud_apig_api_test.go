package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getApiFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return apis.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccApi_basic(t *testing.T) {
	var (
		api apis.APIResp

		rName       = "huaweicloud_apig_api.test"
		name        = acceptance.RandomAccResourceName()
		updateName  = acceptance.RandomAccResourceName()
		basicConfig = testAccApi_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&api,
		getApiFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApi_basic(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "Public"),
					resource.TestCheckResourceAttr(rName, "description", "Created by script"),
					resource.TestCheckResourceAttr(rName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "request_path", "/user_info/{user_age}"),
					resource.TestCheckResourceAttr(rName, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(rName, "matching", "Exact"),
					resource.TestCheckResourceAttr(rName, "success_response", "Success response"),
					resource.TestCheckResourceAttr(rName, "failure_response", "Failed response"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "backend_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "web.0.path", "/getUserAge/{userAge}"),
					resource.TestCheckResourceAttr(rName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "web.0.timeout", "30000"),
					resource.TestCheckResourceAttr(rName, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(rName, "web_policy.0.conditions.#", "1"),
					resource.TestCheckResourceAttr(rName, "mock.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(rName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph_policy.#", "0"),
					resource.TestCheckResourceAttrPair(rName, "web.0.authorizer_id",
						"huaweicloud_apig_custom_authorizer.test", "id"),
					resource.TestCheckResourceAttr(rName, "web_policy.0.backend_params.#", "3"),
				),
			},
			{
				Config: testAccApi_update(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "type", "Public"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(rName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "request_path", "/user_info/{user_name}"),
					resource.TestCheckResourceAttr(rName, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(rName, "matching", "Exact"),
					resource.TestCheckResourceAttr(rName, "success_response", "Updated Success response"),
					resource.TestCheckResourceAttr(rName, "failure_response", "Updated Failed response"),
					resource.TestCheckResourceAttr(rName, "request_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "backend_params.#", "3"),
					resource.TestCheckResourceAttr(rName, "web.0.path", "/getUserName/{userName}"),
					resource.TestCheckResourceAttr(rName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "web.0.timeout", "60000"),
					resource.TestCheckResourceAttr(rName, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(rName, "web_policy.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "mock.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(rName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "web_policy.0.backend_params.#", "3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApiResourceImportStateFunc(),
			},
		},
	})
}

func testAccApiResourceImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_api.test"
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("missing some attributes, want '{instance_id}/{name}', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApi_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = "%[2]s_${count.index}"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  availability_zones    = try(slice(data.huaweicloud_availability_zones.test.names, 0, 1), null)
}

resource "huaweicloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
}

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_apig_instance.test.id
  port        = 80
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  dynamic "members" {
    for_each = huaweicloud_compute_instance.test[*].id

    content {
      id = members.value
    }
  }
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[2]s"
  app         = "default"
  description = "API custom authorization test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python3.6"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler(event, context):
    if event["headers"]["authorization"]=='Basic dXNlcjE6cGFzc3dvcmQ=':
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"allow",
                "context":{
                    "user_name":"user1"
                }
            })
        }
    else:
        return {
            'statusCode': 200,
            'body': json.dumps({
                "status":"deny",
                "context":{
                    "code":"1001",
                    "message":"incorrect username or password"
                }
            })
        }
EOF
}

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = huaweicloud_apig_instance.test.id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "BACKEND"
  cache_age        = 60
}
`, common.TestBaseComputeResources(name), name)
}

func testAccApi_basic(relatedConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
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
  request_params {
    name        = "X-TEST-ENUM"
    type        = "STRING"
    location    = "HEADER"
    maximum     = 20
    minimum     = 10
    example     = "ACC_TEST_XXX"
    passthrough = true
    enumeration = "ACC_TEST_A,ACC_TEST_B"
  }

  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }
  backend_params {
    type              = "SYSTEM"
    name              = "x-test-id"
    location          = "HEADER"
    value             = "x-test-id"
    system_param_type = "backend"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id
  }

  web_policy {
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getUserAge/{userAge}"
    timeout          = 30000
    retry_count      = 1
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "x-test-policy-id"
      location          = "HEADER"
      value             = "x-test-policy-id"
      system_param_type = "backend"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "%[2]s"
      location          = "HEADER"
      value             = "serverName"
      system_param_type = "internal"
    }

    conditions {
      source     = "param"
      param_name = "user_age"
      type       = "Equal"
      value      = "28"
    }
  }
}
`, relatedConfig, name)
}

func testAccApi_update(relatedConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/{user_name}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Updated Success response"
  failure_response        = "Updated Failed response"
  description             = "Updated by script"

  request_params {
    name     = "user_name"
    type     = "STRING"
    location = "PATH"
    required = true
    maximum  = 64
    minimum  = 3
  }
  request_params {
    name        = "X-TEST-ENUM"
    type        = "STRING"
    location    = "HEADER"
    maximum     = 20
    minimum     = 10
    example     = "ACC_TEST_XXXX"
    passthrough = false
    enumeration = "ACC_TEST_A,ACC_TEST_B,ACC_TEST_C"
  }

  backend_params {
    type     = "REQUEST"
    name     = "userName"
    location = "PATH"
    value    = "user_name"
  }
  backend_params {
    type              = "SYSTEM"
    name              = "x-update-policy-id"
    location          = "HEADER"
    value             = "x-update-policy-id"
    system_param_type = "backend"
  }
  backend_params {
    type              = "SYSTEM"
    name              = "%[2]s"
    location          = "HEADER"
    value             = "serverName"
    system_param_type = "internal"
  }

  web {
    path             = "/getUserName/{userName}"
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 60000
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id
  }

  web_policy {
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getAdminName/{adminName}"
    timeout          = 60000
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id
    authorizer_id    = huaweicloud_apig_custom_authorizer.test.id

    backend_params {
      type     = "REQUEST"
      name     = "adminName"
      location = "PATH"
      value    = "user_name"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "x-update-policy-id"
      location          = "HEADER"
      value             = "x-update-policy-id"
      system_param_type = "backend"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "%[2]s"
      location          = "HEADER"
      value             = "serverName"
      system_param_type = "internal"
    }

    conditions {
      source     = "param"
      param_name = "user_name"
      type       = "Equal"
      value      = "Administrator"
    }
    conditions {
      source      = "cookie"
      cookie_name = "user_name"
      type        = "Equal"
      value       = "value_test"
    }
  }
}
`, relatedConfig, name)
}
