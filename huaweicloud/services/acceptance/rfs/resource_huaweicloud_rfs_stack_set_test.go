package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rfs"
)

func getStackSetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region       = acceptance.HW_REGION_NAME
		stackSetName = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	return rfs.QueryStackSetMetaData(client, stackSetName)
}

func TestAccStackSet_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_rfs_stack_set.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStackSetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testStackSet_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "stack_set_name", name),
					resource.TestCheckResourceAttr(rName, "stack_set_description", "stack set description"),
					resource.TestCheckResourceAttr(rName, "permission_model", "SELF_MANAGED"),
					resource.TestCheckResourceAttr(rName, "initial_stack_description", "initial stack description"),
					resource.TestCheckResourceAttr(rName, "managed_operation.0.enable_parallel_operation", "false"),
					resource.TestCheckResourceAttr(rName, "administration_agency_name", "for-rfs-test"),
					resource.TestCheckResourceAttr(rName, "managed_agency_name", "Target_Account-test"),
					resource.TestCheckResourceAttrSet(rName, "stack_set_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testStackSet_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "stack_set_name", name),
					resource.TestCheckResourceAttr(rName, "permission_model", "SELF_MANAGED"),
					resource.TestCheckResourceAttrSet(rName, "stack_set_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttr(rName, "initial_stack_description", "initial stack description update"),
					resource.TestCheckResourceAttr(rName, "stack_set_description", "stack set description update"),
					resource.TestCheckResourceAttr(rName, "administration_agency_name", "for-rfs-test_update"),
					resource.TestCheckResourceAttr(rName, "managed_agency_name", "Target_Account-test_update"),
					resource.TestCheckResourceAttr(rName, "managed_operation.0.enable_parallel_operation", "true"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template_uri", "template_body", "vars_uri", "call_identity"},
			},
		},
	})
}

func testStackSet_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack_set" "test" {
  stack_set_name             = "%[1]s"
  stack_set_description      = "stack set description"
  permission_model           = "SELF_MANAGED"
  template_body              = <<EOT
variable "vpc_name" {
  type    = string
  default = "my-vpc"
}

variable "subnet_name" {
  type    = string
  default = "my-subnet"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet" {
  name       = var.subnet_name
  vpc_id     = huaweicloud_vpc.vpc.id
  cidr       = "172.16.1.0/24"
  gateway_ip = "172.16.1.1"
}
EOT
  administration_agency_name = "for-rfs-test"
  managed_agency_name        = "Target_Account-test"
  initial_stack_description  = "initial stack description"

  managed_operation {
    enable_parallel_operation = false
  }
}
`, name)
}

func testStackSet_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_stack_set" "test" {
  stack_set_name             = "%[1]s"
  stack_set_description      = "stack set description update"
  permission_model           = "SELF_MANAGED"
  template_body              = <<EOT
variable "vpc_name" {
  type    = string
  default = "my-vpc"
}

variable "subnet_name" {
  type    = string
  default = "my-subnet"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet" {
  name       = var.subnet_name
  vpc_id     = huaweicloud_vpc.vpc.id
  cidr       = "172.16.1.0/24"
  gateway_ip = "172.16.1.1"
}
EOT
  administration_agency_name = "for-rfs-test_update"
  managed_agency_name        = "Target_Account-test_update"
  initial_stack_description  = "initial stack description update"

  managed_operation {
    enable_parallel_operation = true
  }
}
`, name)
}
