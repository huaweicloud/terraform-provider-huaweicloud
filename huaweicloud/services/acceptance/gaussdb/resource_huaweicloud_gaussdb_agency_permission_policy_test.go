package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/gaussdb"
)

func getResourceAgencyPermissionPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("opengauss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	return gaussdb.GetAgencyPermissionPolicyInfo(client, state.Primary.ID)
}

func TestAccResourceAgencyPermissionPolicy_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_gaussdb_agency_permission_policy.test"
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceAgencyPermissionPolicyFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAgencyPermissionPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "agency_name", "RDSAccessProjectResource"),
					resource.TestCheckResourceAttr(rName, "bind_role_names.#", "2"),
					resource.TestCheckResourceAttr(rName, "unbind_role_names.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "roles.#"),
					resource.TestCheckResourceAttrSet(rName, "is_existed"),
				),
			},
			{
				Config: testAccAgencyPermissionPolicy_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bind_role_names.#", "1"),
					resource.TestCheckResourceAttr(rName, "unbind_role_names.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "roles.#"),
					resource.TestCheckResourceAttrSet(rName, "is_existed"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bind_role_names", "unbind_role_names"},
			},
		},
	})
}

func testAccAgencyPermissionPolicy_basic() string {
	return `
resource "huaweicloud_gaussdb_agency_permission_policy" "test" {
  agency_name       = "RDSAccessProjectResource"
  bind_role_names   = ["DBS AgencyPolicy","GaussDB FullAccess"]
  unbind_role_names = ["GaussDB ReadOnlyAccess"]
}
`
}

func testAccAgencyPermissionPolicy_update() string {
	return `
resource "huaweicloud_gaussdb_agency_permission_policy" "test" {
  agency_name       = "RDSAccessProjectResource"
  bind_role_names   = ["DBS AgencyPolicy"]
  unbind_role_names = ["GaussDB ReadOnlyAccess","GaussDB FullAccess"]
}
`
}
