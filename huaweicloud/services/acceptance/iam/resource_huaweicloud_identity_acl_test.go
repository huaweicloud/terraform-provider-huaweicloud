package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/acl"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAclResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	switch state.Primary.Attributes["type"] {
	case "console":
		v, err := acl.ConsoleACLPolicyGet(client, state.Primary.ID).ConsoleExtract()
		if err != nil {
			return nil, err
		}
		if len(v.AllowAddressNetmasks) == 0 && len(v.AllowIPRanges) == 1 &&
			v.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
			return nil, fmt.Errorf("identity ACL for console access <%s> not exists", state.Primary.ID)
		}
		return v, nil
	case "api":
		v, err := acl.APIACLPolicyGet(client, state.Primary.ID).APIExtract()
		if err != nil {
			return nil, err
		}
		if len(v.AllowAddressNetmasks) == 0 && len(v.AllowIPRanges) == 1 &&
			v.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
			return nil, fmt.Errorf("identity ACL for console access <%s> not exists", state.Primary.ID)
		}
		return v, nil
	}
	return nil, nil
}

func TestAccAcl_basic(t *testing.T) {
	var (
		obj interface{}

		aclByConsole = "huaweicloud_identity_acl.test.0"
		rcByConsole  = acceptance.InitResourceCheck(aclByConsole, &obj, getAclResourceFunc)

		aclByApi = "huaweicloud_identity_acl.test.1"
		rcByApi  = acceptance.InitResourceCheck(aclByApi, &obj, getAclResourceFunc)
	)

	// the runner public IP must by set
	// otherwise, when the ACL is applied, you can't access your account
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckRunnerPublicIP(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcByConsole.CheckResourceDestroy(),
			rcByApi.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAcl_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rcByConsole.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByConsole, "type", "console"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.#", "1"),
					rcByApi.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByApi, "type", "api"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccAcl_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rcByConsole.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByConsole, "type", "console"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.#", "1"),
					rcByApi.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByApi, "type", "api"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.#", "1"),
				),
			},
		},
	})
}

func testAccAcl_basic_step1() string {
	return fmt.Sprintf(`
variable "acl_types" {
  type    = list(string)
  default = ["console", "api"]
}

resource "huaweicloud_identity_acl" "test" {
  count = 2

  type = var.acl_types[count.index]

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a basic ip range for ${var.acl_types[count.index]} access"
  }

  ip_cidrs {
    cidr        = "%[1]s/32"
    description = "This is a basic ip address for ${var.acl_types[count.index]} access"
  }
}
`, acceptance.HW_RUNNER_PUBLIC_IP)
}

func testAccAcl_basic_step2() string {
	return fmt.Sprintf(`
variable "acl_types" {
  type    = list(string)
  default = ["console", "api"]
}

resource "huaweicloud_identity_acl" "test" {
  count = 2

  type = var.acl_types[count.index]

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a update ip range 1 for ${var.acl_types[count.index]} access"
  }
  ip_ranges {
    range       = "192.168.0.0-192.168.255.255"
    description = "This is a update ip range 2 for ${var.acl_types[count.index]} access"
  }

  ip_cidrs {
    cidr        = "%[1]s/32"
    description = "This is a update ip address for ${var.acl_types[count.index]} access"
  }
}
`, acceptance.HW_RUNNER_PUBLIC_IP)
}
