package iam

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getAclResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetAclByDomainId(client, state.Primary.Attributes["type"], cfg.DomainID)
}

func TestAccAcl_basic(t *testing.T) {
	var (
		obj interface{}

		aclByConsole = "huaweicloud_identity_acl.test.0"
		rcByConsole  = acceptance.InitResourceCheck(aclByConsole, &obj, getAclResourceFunc)

		aclByApi = "huaweicloud_identity_acl.test.1"
		rcByApi  = acceptance.InitResourceCheck(aclByApi, &obj, getAclResourceFunc)
	)

	// the runner public IP must includes the IP address of your current machine and make sure the first IP address is
	// the IP address of your current operating machine
	// Otherwise, when the ACL is applied, you can't access with your current operating machine
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckRunnerPublicIPs(t, 2)
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
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.#", "3"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.0.range", "172.16.0.0-172.16.0.255"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.0.description",
						"This is a basic IP range for 172.16.0.0/24, which in console access"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.1.range", "192.168.0.0-192.168.0.255"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.1.description",
						"This is a basic IP range for 192.168.0.0/24, which in console access"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.2.range", "10.16.0.0-10.16.0.255"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.2.description",
						"This is a basic IP range for 10.16.0.0/24, which in console access"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.#", "1"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.0.cidr",
						fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0])),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.0.description",
						fmt.Sprintf("This is a basic IPv4 CIDR block %s for console access",
							fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0]))),
					rcByApi.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByApi, "type", "api"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.#", "3"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.0.range", "172.16.0.0-172.16.0.255"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.0.description",
						"This is a basic IP range for 172.16.0.0/24, which in api access"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.1.range", "192.168.0.0-192.168.0.255"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.1.description",
						"This is a basic IP range for 192.168.0.0/24, which in api access"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.2.range", "10.16.0.0-10.16.0.255"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.2.description",
						"This is a basic IP range for 10.16.0.0/24, which in api access"),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.#", "1"),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.0.cidr",
						fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0])),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.0.description",
						fmt.Sprintf("This is a basic IPv4 CIDR block %s for api access",
							fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0]))),
				),
			},
			{
				Config: testAccAcl_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rcByConsole.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByConsole, "type", "console"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.#", "3"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.0.range", "172.16.0.0-172.16.255.255"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.0.description",
						"This is a updated IP range for 172.16.0.0/16, which in console access"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.1.range", "192.168.0.0-192.168.255.255"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.1.description",
						"This is a updated IP range for 192.168.0.0/16, which in console access"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.2.range", "10.16.0.0-10.16.255.255"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_ranges.2.description",
						"This is a updated IP range for 10.16.0.0/16, which in console access"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.#", "2"),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.0.cidr",
						fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0])),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.0.description",
						fmt.Sprintf("This is a basic IPv4 CIDR block %s for console access",
							fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0]))),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.1.cidr",
						fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[1])),
					resource.TestCheckResourceAttr(aclByConsole, "ip_cidrs.1.description",
						fmt.Sprintf("This is a basic IPv4 CIDR block %s for console access",
							fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[1]))),
					rcByApi.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclByApi, "type", "api"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.#", "3"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.0.range", "172.16.0.0-172.16.255.255"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.0.description",
						"This is a updated IP range for 172.16.0.0/16, which in api access"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.1.range", "192.168.0.0-192.168.255.255"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.1.description",
						"This is a updated IP range for 192.168.0.0/16, which in api access"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.2.range", "10.16.0.0-10.16.255.255"),
					resource.TestCheckResourceAttr(aclByApi, "ip_ranges.2.description",
						"This is a updated IP range for 10.16.0.0/16, which in api access"),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.#", "2"),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.0.cidr",
						fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0])),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.0.description",
						fmt.Sprintf("This is a basic IPv4 CIDR block %s for api access",
							fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[0]))),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.1.cidr",
						fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[1])),
					resource.TestCheckResourceAttr(aclByApi, "ip_cidrs.1.description",
						fmt.Sprintf("This is a basic IPv4 CIDR block %s for api access",
							fmt.Sprintf("%s/32", strings.Split(acceptance.HW_RUNNER_PUBLIC_IPS, ",")[1]))),
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

variable "runner_public_ips" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_identity_acl" "test" {
  count = 2

  type = var.acl_types[count.index]

  ip_ranges {
    range       = "172.16.0.0-172.16.0.255"
    description = "This is a basic IP range for 172.16.0.0/24, which in ${var.acl_types[count.index]} access"
  }
  ip_ranges {
    range       = "192.168.0.0-192.168.0.255"
    description = "This is a basic IP range for 192.168.0.0/24, which in ${var.acl_types[count.index]} access"
  }
  ip_ranges {
	range       = "10.16.0.0-10.16.0.255"
	description = "This is a basic IP range for 10.16.0.0/24, which in ${var.acl_types[count.index]} access"
  }

  dynamic "ip_cidrs" {
    for_each = var.runner_public_ips != "" ? slice(split(",", var.runner_public_ips), 0, 1) : []

    content {
      cidr        = format("%%s/32", ip_cidrs.value)
      description = "This is a basic IPv4 CIDR block ${format("%%s/32", ip_cidrs.value)} for ${var.acl_types[count.index]} access"
    }
  }
}
`, acceptance.HW_RUNNER_PUBLIC_IPS)
}

func testAccAcl_basic_step2() string {
	return fmt.Sprintf(`
variable "acl_types" {
  type    = list(string)
  default = ["console", "api"]
}

variable "runner_public_ips" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_identity_acl" "test" {                                                                                              
  count = 2

  type = var.acl_types[count.index]

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a updated IP range for 172.16.0.0/16, which in ${var.acl_types[count.index]} access"
  }
  ip_ranges {
    range       = "192.168.0.0-192.168.255.255"
    description = "This is a updated IP range for 192.168.0.0/16, which in ${var.acl_types[count.index]} access"
  }
  ip_ranges {
	range       = "10.16.0.0-10.16.255.255"
	description = "This is a updated IP range for 10.16.0.0/16, which in ${var.acl_types[count.index]} access"
  }

  dynamic "ip_cidrs" {
    for_each = var.runner_public_ips != "" ? slice(split(",", var.runner_public_ips), 0, 2) : []

    content {
      cidr        = format("%%s/32", ip_cidrs.value)
      description = "This is a basic IPv4 CIDR block ${format("%%s/32", ip_cidrs.value)} for ${var.acl_types[count.index]} access"
    }
  }
}
`, acceptance.HW_RUNNER_PUBLIC_IPS)
}
