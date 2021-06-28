package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3.0/acl"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIdentitACL_basic(t *testing.T) {
	var acl acl.ACLPolicy
	resourceName := "huaweicloud_identity_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityACL_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityACLExists(resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "type", "console"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIdentityACL_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityACLExists(resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "type", "console"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "2"),
				),
			},
		},
	})
}

func TestAccIdentitACL_apiAccess(t *testing.T) {
	var acl acl.ACLPolicy
	resourceName := "huaweicloud_identity_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityACL_apiAccess(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityACLExists(resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "type", "api"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIdentityACL_apiUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityACLExists(resourceName, &acl),
					resource.TestCheckResourceAttr(resourceName, "type", "api"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIdentityACLExists(n string, ac *acl.ACLPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.IAMV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud IAM client: %s", err)
		}

		switch rs.Primary.Attributes["type"] {
		case "console":
			v, err := acl.ConsoleACLPolicyGet(client, rs.Primary.ID).ConsoleExtract()
			if err != nil {
				return err
			}
			if len(v.AllowAddressNetmasks) == 0 && len(v.AllowIPRanges) == 1 &&
				v.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
				return fmtp.Errorf("Identity ACL for console access <%s> not exists", rs.Primary.ID)
			}
			ac = v
		case "api":
			v, err := acl.APIACLPolicyGet(client, rs.Primary.ID).APIExtract()
			if err != nil {
				return err
			}
			if len(v.AllowAddressNetmasks) == 0 && len(v.AllowIPRanges) == 1 &&
				v.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
				return fmtp.Errorf("Identity ACL for console access <%s> not exists", rs.Primary.ID)
			}
			ac = v
		}

		return nil
	}
}

func testAccCheckIdentityACLDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.IAMV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IAM client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_identity_acl" {
			continue
		}
		switch rs.Primary.Attributes["type"] {
		case "console":
			v, err := acl.ConsoleACLPolicyGet(client, rs.Primary.ID).ConsoleExtract()
			if err == nil && len(v.AllowAddressNetmasks) == len(rs.Primary.Attributes["ip_cidrs.#"]) &&
				(len(v.AllowIPRanges) > 1 || (len(v.AllowIPRanges) == 1 &&
					v.AllowIPRanges[0].IPRange != "0.0.0.0-255.255.255.255")) {
				return fmtp.Errorf("Identity ACL for console access <%s> still exists", rs.Primary.ID)
			}
		case "api":
			v, err := acl.APIACLPolicyGet(client, rs.Primary.ID).APIExtract()
			if err == nil && len(v.AllowAddressNetmasks) == len(rs.Primary.Attributes["ip_cidrs.#"]) &&
				(len(v.AllowIPRanges) > 1 || (len(v.AllowIPRanges) == 1 &&
					v.AllowIPRanges[0].IPRange != "0.0.0.0-255.255.255.255")) {
				return fmtp.Errorf("Identity ACL for api access <%s> still exists", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccIdentityACL_basic() string {
	return fmtp.Sprintf(`
resource "huaweicloud_identity_acl" "test" {
  type = "console"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a basic ip range 1 for console access"
  }

  ip_cidrs {
    cidr        = "159.138.32.195/32"
    description = "This is a basic ip address 1 for console access"
  }
}
`)
}

func testAccIdentityACL_update() string {
	return fmtp.Sprintf(`
resource "huaweicloud_identity_acl" "test" {
  type = "console"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a update ip range 1 for console access"
  }
  ip_ranges {
    range       = "192.168.0.0-192.168.255.255"
    description = "This is a update ip range 2 for console access"
  }

  ip_cidrs {
    cidr        = "159.138.32.195/32"
    description = "This is a update ip address 1 for console access"
  }
  ip_cidrs {
    cidr        = "159.138.32.196/32"
    description = "This is a update ip address 2 for console access"
  }
}
`)
}

func testAccIdentityACL_apiAccess() string {
	return fmtp.Sprintf(`
resource "huaweicloud_identity_acl" "test" {
  type = "api"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a basic ip range 1 for api access"
  }

  ip_cidrs {
    cidr        = "159.138.32.195/32"
    description = "This is a basic ip address 1 for api access"
  }
}
`)
}

func testAccIdentityACL_apiUpdate() string {
	return fmtp.Sprintf(`
resource "huaweicloud_identity_acl" "test" {
  type = "api"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a update ip range 1 for api access"
  }
  ip_ranges {
    range       = "192.168.0.0-192.168.255.255"
    description = "This is a update ip range 2 for api access"
  }

  ip_cidrs {
    cidr        = "159.138.32.195/32"
    description = "This is a update ip address 1 for api access"
  }
  ip_cidrs {
    cidr        = "159.138.32.196/32"
    description = "This is a update ip address 2 for api access"
  }
}
`)
}
