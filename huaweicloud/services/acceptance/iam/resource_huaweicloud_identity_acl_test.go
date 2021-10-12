package iam

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/acl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getIdentitACLResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud IAM client: %s", err)
	}

	switch state.Primary.Attributes["type"] {
	case "console":
		v, err := acl.ConsoleACLPolicyGet(client, state.Primary.ID).ConsoleExtract()
		if err != nil {
			return nil, err
		}
		if len(v.AllowAddressNetmasks) == 0 && len(v.AllowIPRanges) == 1 &&
			v.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
			return nil, fmtp.Errorf("Identity ACL for console access <%s> not exists", state.Primary.ID)
		}
		return v, nil
	case "api":
		v, err := acl.APIACLPolicyGet(client, state.Primary.ID).APIExtract()
		if err != nil {
			return nil, err
		}
		if len(v.AllowAddressNetmasks) == 0 && len(v.AllowIPRanges) == 1 &&
			v.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
			return nil, fmtp.Errorf("Identity ACL for console access <%s> not exists", state.Primary.ID)
		}
		return v, nil
	}
	return nil, nil
}

func TestAccIdentitACL_basic(t *testing.T) {
	var acl acl.ACLPolicy
	resourceName := "huaweicloud_identity_acl.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&acl,
		getIdentitACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityACL_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "console"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIdentityACL_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

	rc := acceptance.InitResourceCheck(
		resourceName,
		&acl,
		getIdentitACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityACL_apiAccess(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "api"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIdentityACL_apiUpdate(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "api"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "2"),
				),
			},
		},
	})
}

func testAccIdentityACL_basic() string {
	return fmt.Sprintf(`
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
	return fmt.Sprintf(`
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
	return fmt.Sprintf(`
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
	return fmt.Sprintf(`
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
