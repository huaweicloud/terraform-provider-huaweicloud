package fgs

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FgsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud FunctionGraph client: %s", err)
	}
	return function.GetMetadata(c, state.Primary.ID).Extract()
}

func TestAccFgsV2Function_basic(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "functiongraph_version", "v1"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				Config: testAccFgsV2Function_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "fuction test update"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_withEpsId(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_withEpsId(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_text(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_text(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"func_code",
				},
			},
		},
	})
}

func TestAccFgsV2Function_createByImage(t *testing.T) {
	var f function.Function
	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_fgs_function.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&f,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckComponentDeployment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_createByImage(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "agency", "functiongraph_swr_trust"),
					resource.TestCheckResourceAttr(resourceName, "runtime", "Custom Image"),
					resource.TestCheckResourceAttr(resourceName, "custom_image.0.url", acceptance.HW_BUILD_IMAGE_URL),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"app",
					"package",
					"xrole",
					"agency",
				},
			},
		},
	})
}

func testAccFgsV2Function_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName)
}

func testAccFgsV2Function_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test update"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName)
}

func testAccFgsV2Function_text(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  description = "fuction test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"

  func_code = <<EOF
# -*- coding:utf-8 -*-
import json
def handler (event, context):
    return {
        "statusCode": 200,
        "isBase64Encoded": False,
        "body": json.dumps(event),
        "headers": {
            "Content-Type": "application/json"
        }
    }
EOF
}
`, rName)
}

func testAccFgsV2Function_withEpsId(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name                  = "%s"
  app                   = "default"
  description           = "fuction test"
  handler               = "index.handler"
  memory_size           = 128
  timeout               = 3
  runtime               = "Python2.7"
  code_type             = "inline"
  enterprise_project_id = "%s"
  func_code             = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccFgsV2Function_createByImage(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  app         = "default"
  memory_size = 128
  runtime     = "Custom Image"
  timeout     = 3
  agency      = "functiongraph_swr_trust"

  custom_image {
    url = "%[2]s"
  }
}
`, rName, acceptance.HW_BUILD_IMAGE_URL)
}
