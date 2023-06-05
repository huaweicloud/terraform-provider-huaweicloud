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

func getOrganizationsUnitResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getOrganizationsUnit: Query Organizations unit
	var (
		getOrganizationsUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		getOrganizationsUnitProduct = "organizations"
	)
	getOrganizationsUnitClient, err := cfg.NewServiceClient(getOrganizationsUnitProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationsUnitPath := getOrganizationsUnitClient.Endpoint + getOrganizationsUnitHttpUrl
	getOrganizationsUnitPath = strings.ReplaceAll(getOrganizationsUnitPath, "{organizational_unit_id}",
		state.Primary.ID)

	getOrganizationsUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationsUnitResp, err := getOrganizationsUnitClient.Request("GET",
		getOrganizationsUnitPath, &getOrganizationsUnitOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations: %s", err)
	}
	return utils.FlattenResponse(getOrganizationsUnitResp)
}

func TestAccOrganizationsUnit_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_organizations_unit.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrganizationsUnitResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOrganizationsUnit_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_organizations.test", "root_id"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
			{
				Config: testOrganizationsUnit_basic_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_organizations.test", "root_id"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
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

func testOrganizationsUnit_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_organizations_unit" "test" {
  name      = "%s"
  parent_id = huaweicloud_organizations.test.root_id

  tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}
`, testOrganizations_basic(), name)
}

func testOrganizationsUnit_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_organizations_unit" "test" {
  name      = "%s"
  parent_id = huaweicloud_organizations.test.root_id

  tags = {
    "key3" = "value3"
    "key4" = "value4"
  }
}
`, testOrganizations_basic(), name)
}
