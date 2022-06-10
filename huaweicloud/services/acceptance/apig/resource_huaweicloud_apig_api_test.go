package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/apis"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigAPIV2_basic(t *testing.T) {
	var (
		// The dedicated instance name only allow letters, digits and underscores (_).
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_api.test"
		api          apis.APIResp
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigAPIDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigAPI_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigAPIExists(resourceName, &api),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "Public"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "request_path", "/user_info/{user_age}"),
					resource.TestCheckResourceAttr(resourceName, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(resourceName, "matching", "Exact"),
					resource.TestCheckResourceAttr(resourceName, "success_response", "Success response"),
					resource.TestCheckResourceAttr(resourceName, "failure_response", "Failed response"),
					resource.TestCheckResourceAttr(resourceName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web.0.path", "/getUserAge/{userAge}"),
					resource.TestCheckResourceAttr(resourceName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "web.0.timeout", "30000"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "mock"),
					resource.TestCheckNoResourceAttr(resourceName, "func_graph"),
					resource.TestCheckNoResourceAttr(resourceName, "mock_policy"),
					resource.TestCheckNoResourceAttr(resourceName, "func_graph_policy"),
				),
			},
			{
				Config: testAccApigAPI_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigAPIExists(resourceName, &api),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "type", "Public"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(resourceName, "request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "request_path", "/user_info/{user_name}"),
					resource.TestCheckResourceAttr(resourceName, "security_authentication", "APP"),
					resource.TestCheckResourceAttr(resourceName, "matching", "Exact"),
					resource.TestCheckResourceAttr(resourceName, "success_response", "Updated Success response"),
					resource.TestCheckResourceAttr(resourceName, "failure_response", "Updated Failed response"),
					resource.TestCheckResourceAttr(resourceName, "request_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "backend_params.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "web.0.path", "/getUserName/{userName}"),
					resource.TestCheckResourceAttr(resourceName, "web.0.request_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "web.0.request_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "web.0.timeout", "60000"),
					resource.TestCheckResourceAttr(resourceName, "web_policy.#", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "mock"),
					resource.TestCheckNoResourceAttr(resourceName, "func_graph"),
					resource.TestCheckNoResourceAttr(resourceName, "mock_policy"),
					resource.TestCheckNoResourceAttr(resourceName, "func_graph_policy"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigAPIResourceImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigAPIDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_api" {
			continue
		}
		_, err := apis.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigAPIExists(n string, app *apis.APIResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Unable to find API ID")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := apis.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*app = *found
		return nil
	}
}

func testAccApigAPIResourceImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", rs.Primary.Attributes["instance_id"],
				rs.Primary.Attributes["name"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["name"]), nil
	}
}

func testAccApigAPI_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_group" "test" {
  name        = "%s"
  instance_id = huaweicloud_apig_instance.test.id
  description = "Created by script"
}

resource "huaweicloud_apig_vpc_channel" "test" {
  name        = "%s"
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
`, testAccApigVpcChannel_base(rName), rName, rName)
}

func testAccApigAPI_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%s"
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
    name             = "%s_policy1"
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
`, testAccApigAPI_base(rName), rName, rName)
}

func testAccApigAPI_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_api" "test" {
  instance_id             = huaweicloud_apig_instance.test.id
  group_id                = huaweicloud_apig_group.test.id
  name                    = "%s_update"
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
    name             = "%s_update_policy1"
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
`, testAccApigAPI_base(rName), rName, rName)
}
