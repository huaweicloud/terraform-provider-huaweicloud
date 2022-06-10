package apig

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/authorizers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigCustomAuthorizerV2_basic(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the authorizer name, environment name
		// and dedicated instance name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_custom_authorizer.test"
		auth         authorizers.CustomAuthorizer
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigCustomAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigCustomAuthorizer_front(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigCustomAuthorizerExists(resourceName, &auth),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "FRONTEND"),
					resource.TestCheckResourceAttr(resourceName, "is_body_send", "true"),
					resource.TestCheckResourceAttr(resourceName, "cache_age", "60"),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "1"),
				),
			},
			{
				Config: testAccApigCustomAuthorizer_frontUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigCustomAuthorizerExists(resourceName, &auth),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "type", "FRONTEND"),
					resource.TestCheckResourceAttr(resourceName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(resourceName, "cache_age", "0"),
					resource.TestCheckResourceAttr(resourceName, "identity.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func TestAccApigCustomAuthorizerV2_backend(t *testing.T) {
	var (
		// Only letters, digits and underscores (_) are allowed in the environment name and dedicated instance name.
		rName        = fmt.Sprintf("tf_acc_test_%s", acctest.RandString(5))
		resourceName = "huaweicloud_apig_custom_authorizer.test"
		auth         authorizers.CustomAuthorizer
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckApigCustomAuthorizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApigCustomAuthorizer_backend(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigCustomAuthorizerExists(resourceName, &auth),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "BACKEND"),
					resource.TestCheckResourceAttr(resourceName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(resourceName, "cache_age", "60"),
				),
			},
			{
				Config: testAccApigCustomAuthorizer_backendUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApigCustomAuthorizerExists(resourceName, &auth),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "type", "BACKEND"),
					resource.TestCheckResourceAttr(resourceName, "is_body_send", "false"),
					resource.TestCheckResourceAttr(resourceName, "cache_age", "45"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccApigSubResNameImportStateFunc(resourceName),
			},
		},
	})
}

func testAccCheckApigCustomAuthorizerDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_apig_custom_authorizer" {
			continue
		}
		_, err := authorizers.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("APIG v2 API custom authorizer (%s) is still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckApigCustomAuthorizerExists(name string, env *authorizers.CustomAuthorizer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No custom authorizer Id")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.ApigV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
		}
		found, err := authorizers.Get(client, rs.Primary.Attributes["instance_id"], rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting custom authorizer (%s): %s", rs.Primary.ID, err)
		}
		*env = *found
		return nil
	}
}

func testAccApigCustomAuthorizer_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
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
`, testAccApigApplication_base(rName), rName)
}

func testAccApigCustomAuthorizer_front(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id  = huaweicloud_apig_instance.test.id
  name         = "%s"
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "FRONTEND"
  is_body_send = true
  cache_age    = 60
  
  identity {
    name     = "user_name"
    location = "QUERY"
  }
}
`, testAccApigCustomAuthorizer_base(rName), rName)
}

func testAccApigCustomAuthorizer_frontUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id  = huaweicloud_apig_instance.test.id
  name         = "%s_update"
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "FRONTEND"
}
`, testAccApigCustomAuthorizer_base(rName), rName)
}

func testAccApigCustomAuthorizer_backend(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id  = huaweicloud_apig_instance.test.id
  name         = "%s"
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "BACKEND"
  cache_age    = 60
}
`, testAccApigCustomAuthorizer_base(rName), rName)
}

func testAccApigCustomAuthorizer_backendUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_apig_custom_authorizer" "test" {
  instance_id  = huaweicloud_apig_instance.test.id
  name         = "%s_update"
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "BACKEND"
  cache_age    = 45
}
`, testAccApigCustomAuthorizer_base(rName), rName)
}
