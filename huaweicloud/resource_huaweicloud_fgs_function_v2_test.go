package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				Config: testAccFgsV2Function_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFgsV2FunctionExists(resourceName, &f),
					resource.TestCheckResourceAttr(resourceName, "description", "fuction test update"),
					resource.TestCheckResourceAttrSet(resourceName, "urn"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
		},
	})
}

func TestAccFgsV2Function_withEpsId(t *testing.T) {
	var f function.Function
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFgsV2FunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFgsV2FunctionExists(resourceName, &f),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
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

func TestAccFgsV2Function_agency(t *testing.T) {
	var f function.Function
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_fgs_function.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFgsV2FunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFgsV2Function_agency(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFgsV2FunctionExists(resourceName, &f),
					resource.TestCheckResourceAttr(resourceName, "agency", rName),
					resource.TestCheckResourceAttr(resourceName, "func_mounts.0.mount_type", "sfs"),
					resource.TestCheckResourceAttr(resourceName, "func_mounts.0.status", "active"),
				),
			},
		},
	})
}

func testAccCheckFgsV2FunctionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	fgsClient, err := config.FgsV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_fgs_function" {
			continue
		}

		_, err := function.GetMetadata(fgsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Function still exists")
		}
	}

	return nil
}

func testAccCheckFgsV2FunctionExists(n string, ft *function.Function) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		fgsClient, err := config.FgsV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
		}

		found, err := function.GetMetadata(fgsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.FuncUrn != rs.Primary.ID {
			return fmtp.Errorf("Function not found")
		}

		*ft = *found

		return nil
	}
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
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccFgsV2Function_agency(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.1.0/24"
  gateway_ip = "192.168.1.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_sfs_file_system" "test" {
  share_proto = "NFS"
  size        = 10
  name        = "%s"
  description = "test sfs for fgs"
}

resource "huaweicloud_identity_agency" "test" {
  name                   = "%s"
  description            = "test agency for fgs"
  delegated_service_name = "op_svc_cff"

  project_role {
    project = "cn-north-4"
    roles = [
      "VPC Administrator",
      "SFS Administrator",
    ]
  }
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%s"
  package     = "default"
  description = "fuction test"
  handler     = "test.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgKGV2ZW50LCBjb250ZXh0KToKICAgIG91dHB1dCA9ICdIZWxsbyBtZXNzYWdlOiAnICsganNvbi5kdW1wcyhldmVudCkKICAgIHJldHVybiBvdXRwdXQ="
  agency      = huaweicloud_identity_agency.test.name
  vpc_id      = huaweicloud_vpc.test.id
  network_id  = huaweicloud_vpc_subnet.test.id

  func_mounts {
    mount_type       = "sfs"
    mount_resource   = huaweicloud_sfs_file_system.test.id
    mount_share_path = huaweicloud_sfs_file_system.test.export_location
    local_mount_path = "/mnt"
  }
}
`, rName, rName, rName, rName, rName)
}
