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

func getIdentitACLResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccIdentitACL_basic(t *testing.T) {
	var object acl.ACLPolicy
	resourceName := "huaweicloud_identity_acl.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentitACLResourceFunc,
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
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityACL_basic("console"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "console"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIdentityACL_update("console"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "console"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "2"),
				),
			},
		},
	})
}

func TestAccIdentitACL_apiAccess(t *testing.T) {
	var object acl.ACLPolicy
	resourceName := "huaweicloud_identity_acl.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentitACLResourceFunc,
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
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityACL_basic("api"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "api"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ip_cidrs.#", "1"),
				),
			},
			{
				Config: testAccIdentityACL_update("api"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "api"),
					resource.TestCheckResourceAttr(resourceName, "ip_ranges.#", "2"),
				),
			},
		},
	})
}

func testAccIdentityACL_basic(aclType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_acl" "test" {
  type = "%[1]s"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a basic ip range for %[1]s access"
  }

  ip_cidrs {
    cidr        = "%[2]s/32"
    description = "This is a basic ip address for %[1]s access"
  }
}
`, aclType, acceptance.HW_CODEARTS_PUBLIC_IP_ADDRESS)
}

func testAccIdentityACL_update(aclType string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_acl" "test" {
  type = "%[1]s"

  ip_ranges {
    range       = "172.16.0.0-172.16.255.255"
    description = "This is a update ip range 1 for %[1]s access"
  }
  ip_ranges {
    range       = "192.168.0.0-192.168.255.255"
    description = "This is a update ip range 2 for %[1]s access"
  }

  ip_cidrs {
    cidr        = "%[2]s/32"
    description = "This is a update ip address for %[1]s access"
  }
}
`, aclType, acceptance.HW_CODEARTS_PUBLIC_IP_ADDRESS)
}
