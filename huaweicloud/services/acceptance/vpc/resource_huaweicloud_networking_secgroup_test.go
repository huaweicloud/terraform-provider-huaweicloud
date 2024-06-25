package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNetworkSecGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud VPC network v1 client: %s", err)
	}

	return securitygroups.Get(client, state.Primary.ID).Extract()
}

func TestAccNetworkingV3SecGroup_basic(t *testing.T) {
	var secGroup securitygroups.SecurityGroup
	name := acceptance.RandomAccResourceNameWithDash()
	updatedName := fmt.Sprintf("%s-updated", name)
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secGroup,
		getNetworkSecGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
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
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
		},
	})
}

func TestAccNetworkingV3SecGroup_withEpsId(t *testing.T) {
	var secGroup securitygroups.SecurityGroup
	name := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secGroup,
		getNetworkSecGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_epsId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccNetworkingV3SecGroup_noDefaultRules(t *testing.T) {
	var secGroup securitygroups.SecurityGroup
	name := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_networking_secgroup.secgroup_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secGroup,
		getNetworkSecGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_noDefaultRules(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "0"),
				),
			},
		},
	})
}

func testAccSecGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "security group acceptance test"

  tags = {
	foo = "bar"
	key = "value"
  }
}
`, name)
}

func testAccSecGroup_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "security group acceptance test updated"

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
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
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
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
