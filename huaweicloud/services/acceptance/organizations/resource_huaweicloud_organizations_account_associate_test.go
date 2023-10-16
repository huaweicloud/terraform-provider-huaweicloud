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

func getAccountAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAccount: Query Organizations account
	var (
		getAccountHttpUrl = "v1/organizations/accounts/{account_id}"
		getAccountProduct = "organizations"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations Client: %s", err)
	}

	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{account_id}", state.Primary.ID)

	getAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountResp, err := getAccountClient.Request("GET", getAccountPath, &getAccountOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AccountAssociate: %s", err)
	}
	return utils.FlattenResponse(getAccountResp)
}

func TestAccAccountAssociate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_organizations_account_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountAssociateResourceFunc,
	)

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsOrganizationalUnitId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccountAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "account_id",
						"huaweicloud_organizations_account.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_organizations_organizational_unit.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
				),
			},
			{
				Config: testAccountAssociate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "account_id",
						"huaweicloud_organizations_account.test", "id"),
					resource.TestCheckResourceAttr(rName, "parent_id",
						acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccountAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_organizations_account_associate" "test" {
  account_id = huaweicloud_organizations_account.test.id
  parent_id  = huaweicloud_organizations_organizational_unit.test.id
}
`, testOrganizationalUnit_basic(name), testAccount_basic(name))
}

func testAccountAssociate_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_account_associate" "test" {
  account_id = huaweicloud_organizations_account.test.id
  parent_id  = "%s"
}
`, testAccount_basic(name), name, acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID)
}
