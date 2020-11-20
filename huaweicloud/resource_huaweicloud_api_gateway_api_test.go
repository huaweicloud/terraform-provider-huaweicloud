package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/apigw/apis"
)

func TestAccApiGatewayAPI_basic(t *testing.T) {
	var resName = "huaweicloud_api_gateway_api.acc_apigw_api"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApiGatewayApiDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigwAPI_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayApiExists(resName),
					resource.TestCheckResourceAttr(resName, "name", "acc_apigw_api"),
					resource.TestCheckResourceAttr(resName, "group_name", "acc_apigw_group_1"),
					resource.TestCheckResourceAttr(resName, "auth_type", "NONE"),
					resource.TestCheckResourceAttr(resName, "backend_type", "HTTP"),
					resource.TestCheckResourceAttr(resName, "request_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resName, "request_method", "GET"),
					resource.TestCheckResourceAttr(resName, "request_uri", "/test/path1"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.method", "GET"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.uri", "/web/openapi"),
					resource.TestCheckResourceAttr(resName, "http_backend.0.timeout", "10000"),
				),
			},
			{
				Config: testAccApigwAPI_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApiGatewayApiExists(resName),
					resource.TestCheckResourceAttr(resName, "description", "updated by acc test"),
					resource.TestCheckResourceAttr(resName, "auth_type", "IAM"),
					resource.TestCheckResourceAttr(resName, "request_protocol", "BOTH"),
					resource.TestCheckResourceAttr(resName, "request_uri", "/test/path2"),
				),
			},
		},
	})
}

func testAccCheckApiGatewayApiDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	apigwClient, err := config.apiGatewayV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_api_gateway_api" {
			continue
		}

		_, err := apis.Get(apigwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("api gateway API still exists")
		}
	}

	return nil
}

func testAccCheckApiGatewayApiExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		apigwClient, err := config.apiGatewayV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud api gateway client: %s", err)
		}

		found, err := apis.Get(apigwClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("apigateway API not found")
		}

		return nil
	}
}

const testAccApigwAPI_basic = `
resource "huaweicloud_api_gateway_group" "acc_apigw_group" {
  name = "acc_apigw_group_1"
  description = "created by acc test"
}

resource "huaweicloud_api_gateway_api" "acc_apigw_api" {
  group_id = huaweicloud_api_gateway_group.acc_apigw_group.id
  name   = "acc_apigw_api"
  description  = "created by acc test"
  tags = ["tag1","tag2"]
  visibility = 2
  auth_type = "NONE"
  backend_type = "HTTP"
  request_protocol = "HTTPS"
  request_method = "GET"
  request_uri = "/test/path1"
  example_success_response = "this is a successful response"

  http_backend {
    protocol = "HTTPS"
    method = "GET"
    uri = "/web/openapi"
    url_domain = "myhuaweicloud.com"
    timeout = 10000
  }
}
`
const testAccApigwAPI_update = `
resource "huaweicloud_api_gateway_group" "acc_apigw_group" {
  name = "acc_apigw_group_1"
  description = "created by acc test"
}

resource "huaweicloud_api_gateway_api" "acc_apigw_api" {
  group_id = huaweicloud_api_gateway_group.acc_apigw_group.id
  name   = "acc_apigw_api"
  description  = "updated by acc test"
  tags = ["tag1","tag2"]
  visibility = 2
  auth_type = "IAM"
  backend_type = "HTTP"
  request_protocol = "BOTH"
  request_method = "GET"
  request_uri = "/test/path2"
  example_success_response = "this is a successful response"

  http_backend {
    protocol = "HTTPS"
    method = "GET"
    uri = "/web/openapi"
    url_domain = "myhuaweicloud.com"
    timeout = 10000
  }
}
`
