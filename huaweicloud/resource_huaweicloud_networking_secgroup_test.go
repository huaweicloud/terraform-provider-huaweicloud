package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingV3SecGroup_basic(t *testing.T) {
	var secGroup securitygroups.SecurityGroup
	name := fmt.Sprintf("seg-acc-test-%s", acctest.RandString(5))
	updatedName := fmt.Sprintf("%s-updated", name)
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV3SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV3SecGroupExists(resourceName, &secGroup),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecGroup_update(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "id", &secGroup.ID),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
		},
	})
}

func TestAccNetworkingV3SecGroup_withEpsId(t *testing.T) {
	var secGroup securitygroups.SecurityGroup
	name := fmt.Sprintf("seg-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV3SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV3SecGroupExists(resourceName, &secGroup),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccNetworkingV3SecGroup_noDefaultRules(t *testing.T) {
	var secGroup securitygroups.SecurityGroup
	name := fmt.Sprintf("seg-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV3SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_noDefaultRules(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV3SecGroupExists(resourceName, &secGroup),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "0"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV3SecGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_secgroup" {
			continue
		}

		_, err := securitygroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Security group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV3SecGroupExists(n string, secGroup *securitygroups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := securitygroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Security group not found")
		}

		*secGroup = *found

		return nil
	}
}

func testAccSecGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "security group acceptance test"
}
`, name)
}

func testAccSecGroup_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "security group acceptance test updated"
}
`, name)
}

func testAccSecGroup_epsId(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name                  = "%s"
  description           = "ecurity group acceptance test with eps ID"
  enterprise_project_id = "%s"
}
`, name, HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccSecGroup_noDefaultRules(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name                 = "%s"
  description          = "security group acceptance test without default rules"
  delete_default_rules = true
}
`, name)
}
