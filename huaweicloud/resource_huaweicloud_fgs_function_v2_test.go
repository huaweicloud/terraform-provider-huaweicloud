package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/fgs/v2/function"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccFgsV2Function_basic(t *testing.T) {
	var f function.Function
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFgsV2FunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFgsV2FunctionExists(resourceName, &f),
				),
			},
		},
	})
}

func TestAccFgsV2Function_text(t *testing.T) {
	var f function.Function
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFgsV2FunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_text(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFgsV2FunctionExists(resourceName, &f),
				),
			},
		},
	})
}

func testAccCheckFgsV2FunctionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	fgsClient, err := config.FgsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_fgs_function" {
			continue
		}

		_, err := function.GetMetadata(fgsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Function still exists")
		}
	}

	return nil
}

func testAccCheckFgsV2FunctionExists(n string, ft *function.Function) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		fgsClient, err := config.FgsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
		}

		found, err := function.GetMetadata(fgsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.FuncUrn != rs.Primary.ID {
			return fmt.Errorf("Function not found")
		}

		*ft = *found

		return nil
	}
}

func testAccFgsV2Function_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name = "%s"
  app = "default"
  description = "fuction test"
  handler = "index.handler"
  memory_size = 128
  timeout = 3
  runtime = "Python2.7"
  code_type = "inline"
  func_code = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
}
`, rName)
}

func testAccFgsV2Function_text(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name = "%s"
  app = "default"
  description = "fuction test"
  handler = "index.handler"
  memory_size = 128
  timeout = 3
  runtime = "Python2.7"
  code_type = "inline"
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
