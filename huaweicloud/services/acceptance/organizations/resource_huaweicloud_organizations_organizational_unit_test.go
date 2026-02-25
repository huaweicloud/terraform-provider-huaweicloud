package organizations

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getOrganizationalUnitResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getOrganizationalUnit: Query Organizations organizational unit
	var (
		region                       = acceptance.HW_REGION_NAME
		getOrganizationalUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		getOrganizationalUnitProduct = "organizations"
	)
	getOrganizationalUnitClient, err := cfg.NewServiceClient(getOrganizationalUnitProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationalUnitPath := getOrganizationalUnitClient.Endpoint + getOrganizationalUnitHttpUrl
	getOrganizationalUnitPath = strings.ReplaceAll(getOrganizationalUnitPath, "{organizational_unit_id}",
		state.Primary.ID)

	getOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationalUnitResp, err := getOrganizationalUnitClient.Request("GET",
		getOrganizationalUnitPath, &getOrganizationalUnitOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations organizational unit: %s", err)
	}
	return utils.FlattenResponse(getOrganizationalUnitResp)
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
