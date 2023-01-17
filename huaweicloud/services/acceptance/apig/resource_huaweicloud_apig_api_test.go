package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
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

		rName      = "huaweicloud_apig_api.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
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
				Config: testAccApi_basic(name),
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
					resource.TestCheckResourceAttr(rName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "web.0.path", "/getUserAge/{userAge}"),
					resource.TestCheckResourceAttr(rName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "web.0.timeout", "30000"),
					resource.TestCheckResourceAttr(rName, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(rName, "mock.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(rName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph_policy.#", "0"),
				),
			},
			{
				Config: testAccApi_update(updateName),
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
					resource.TestCheckResourceAttr(rName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(rName, "web.0.path", "/getUserName/{userName}"),
					resource.TestCheckResourceAttr(rName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(rName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(rName, "web.0.timeout", "60000"),
					resource.TestCheckResourceAttr(rName, "web_policy.#", "1"),
					resource.TestCheckResourceAttr(rName, "mock.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph.#", "0"),
					resource.TestCheckResourceAttr(rName, "mock_policy.#", "0"),
					resource.TestCheckResourceAttr(rName, "func_graph_policy.#", "0"),
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

  members {
    id = huaweicloud_compute_instance.test.id
  }
}
`, testAccVpcChannel_base(name), name)
}

func testAccApi_basic(name string) string {
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
  
  backend_params {
	type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id
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
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id

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
`, testAccApi_base(name), name)
}

func testAccApi_update(name string) string {
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
  
  backend_params {
    type     = "REQUEST"
    name     = "userName"
    location = "PATH"
    value    = "user_name"
  }

  web {
    path             = "/getUserName/{userName}"
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 60000
  }

  web_policy {
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getAdminName/{adminName}"
    timeout          = 60000
    vpc_channel_id   = huaweicloud_apig_vpc_channel.test.id

    backend_params {
      type     = "REQUEST"
      name     = "adminName"
      location = "PATH"
      value    = "user_name"
    }

    conditions {
      source     = "param"
      param_name = "user_name"
      type       = "Equal"
      value      = "Administrator"
    }
  }
}
`, testAccApi_base(name), name)
}
