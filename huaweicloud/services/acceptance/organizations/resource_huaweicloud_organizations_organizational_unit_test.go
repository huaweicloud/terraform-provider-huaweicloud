package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
)

func getOrganizationalUnitResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetOrganizationalUnit(client, state.Primary.ID)
}

func TestAccOrganizationalUnit_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_organizations_organizational_unit.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getOrganizationalUnitResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationalUnit_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testAccOrganizationalUnit_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "tags.%", "1"),
					resource.TestCheckResourceAttr(rName, "tags.foo1", "bar_update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_id"},
			},
		},
	})
}

func testAccOrganizationalUnit_basic_step1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, name)
}

func testAccOrganizationalUnit_basic_step2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id

  tags = {
    foo1 = "bar_update"
  }
}
`, name)
}
