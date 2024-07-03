package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/authorizers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCustomAuthorizerFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return authorizers.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID).Extract()
}

func TestAccCustomAuthorizer_basic(t *testing.T) {
	var (
		auth authorizers.CustomAuthorizer

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		rName      = "huaweicloud_apig_custom_authorizer.test"
		baseConfig = testAccCustomAuthorizer_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&auth,
		getCustomAuthorizerFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomAuthorizer_front_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "FRONTEND"),
					resource.TestCheckResourceAttr(rName, "is_body_send", "true"),
					resource.TestCheckResourceAttr(rName, "cache_age", "60"),
					resource.TestCheckResourceAttr(rName, "identity.#", "1"),
				),
			},
			{
				Config: testAccCustomAuthorizer_front_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "type", "FRONTEND"),
					resource.TestCheckResourceAttr(rName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(rName, "cache_age", "0"),
					resource.TestCheckResourceAttr(rName, "identity.#", "0"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCustomAuthorizerImportStateFunc(),
			},
		},
	})
}

func TestAccCustomAuthorizer_backend(t *testing.T) {
	var (
		auth authorizers.CustomAuthorizer

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
		rName      = "huaweicloud_apig_custom_authorizer.test"
		baseConfig = testAccCustomAuthorizer_base(name)
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&auth,
		getCustomAuthorizerFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomAuthorizer_backend_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "BACKEND"),
					resource.TestCheckResourceAttr(rName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(rName, "cache_age", "60"),
				),
			},
			{
				Config: testAccCustomAuthorizer_backend_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "type", "BACKEND"),
					resource.TestCheckResourceAttr(rName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(rName, "cache_age", "45"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCustomAuthorizerImportStateFunc(),
			},
		},
	})
}

func testAccCustomAuthorizerImportStateFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rName := "huaweicloud_apig_custom_authorizer.test"
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

func testAccCustomAuthorizer_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
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

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}
`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccCustomAuthorizer_front_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "FRONTEND"
  is_body_send     = true
  cache_age        = 60
  
  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
`, baseConfig, name)
}

func testAccCustomAuthorizer_front_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "FRONTEND"
}
`, baseConfig, name)
}

func testAccCustomAuthorizer_backend_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "BACKEND"
  cache_age        = 60
}
`, baseConfig, name)
}

func testAccCustomAuthorizer_backend_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id      = local.instance_id
  name             = "%[2]s"
  function_urn     = huaweicloud_fgs_function.test.urn
  function_version = "latest"
  type             = "BACKEND"
  cache_age        = 45
}
`, baseConfig, name)
}
