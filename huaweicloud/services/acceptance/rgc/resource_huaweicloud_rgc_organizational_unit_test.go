package rgc

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

func getOrganizationalUnitRegisterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getAccount: Query RGC organizational unit via rgc API
	var (
		region            = acceptance.HW_REGION_NAME
		getAccountHttpUrl = "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}"
		getAccountProduct = "rgc"
	)
	getOrganizationalUnitClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RGC client: %s", err)
	}

	getOrganizationalUnitPath := getOrganizationalUnitClient.Endpoint + getAccountHttpUrl
	getOrganizationalUnitPath = strings.ReplaceAll(getOrganizationalUnitPath, "{managed_organizational_unit_id}", state.Primary.ID)

	getOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOrganizationalUnitResp, err := getOrganizationalUnitClient.Request("GET", getOrganizationalUnitPath, &getOrganizationalUnitOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving organizational unit: %s", err)
	}

	return utils.FlattenResponse(getOrganizationalUnitResp)
}

func TestAccOrganizationalUnit_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rgc_organizational_unit.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrganizationalUnitRegisterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOrganizationalUnitID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationalUnit_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "parent_organizational_unit_name"),
					resource.TestCheckResourceAttrSet(rName, "manage_account_id"),
					resource.TestCheckResourceAttrSet(rName, "organizational_unit_id"),
					resource.TestCheckResourceAttrSet(rName, "organizational_unit_status"),
					resource.TestCheckResourceAttrSet(rName, "organizational_unit_type"),
				),
			},
		},
	})
}

func testAccOrganizationalUnit_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_organizational_unit" "test" {
  organizational_unit_name      = "%[1]s"
  parent_organizational_unit_id = "%[2]s"
}
`, name, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
