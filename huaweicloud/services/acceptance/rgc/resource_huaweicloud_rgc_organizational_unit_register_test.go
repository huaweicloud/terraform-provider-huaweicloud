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

func getOrganizationalUnitResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccOrganizationalUnitRegister_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_rgc_organizational_unit_register.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOrganizationalUnitResourceFunc,
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
				Config: testAccOrganizationalUnitRegister_basic(),
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

func testAccOrganizationalUnitRegister_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_organizational_unit_register" "test" {
  organizational_unit_id = "%[1]s"
}
`, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
