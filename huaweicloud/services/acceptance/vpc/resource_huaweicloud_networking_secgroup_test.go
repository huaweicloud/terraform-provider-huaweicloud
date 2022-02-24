package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSecGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return groups.Get(c, state.Primary.ID)
}

func TestAccNetworkingV3SecGroup_basic(t *testing.T) {
	var secGroup groups.SecurityGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_networking_secgroup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secGroup,
		getSecGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSecGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &secGroup.ID),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
				),
			},
		},
	})
}

func TestAccNetworkingV3SecGroup_withEpsId(t *testing.T) {
	var secGroup groups.SecurityGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_networking_secgroup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secGroup,
		getSecGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccNetworkingV3SecGroup_noDefaultRules(t *testing.T) {
	var secGroup groups.SecurityGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_networking_secgroup.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secGroup,
		getSecGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecGroup_noDefaultRules(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "0"),
				),
			},
		},
	})
}

func testAccSecGroup_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "security group acceptance test"
}
`, name)
}

func testAccSecGroup_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s-update"
  description = "security group acceptance test updated"
}
`, name)
}

func testAccSecGroup_epsId(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name                  = "%s"
  description           = "ecurity group acceptance test with eps ID"
  enterprise_project_id = "%s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccSecGroup_noDefaultRules(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%s"
  description          = "security group acceptance test without default rules"
  delete_default_rules = true
}
`, name)
}
