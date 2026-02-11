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

func getAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getAccount: Query Organizations account
	var (
		region            = acceptance.HW_REGION_NAME
		getAccountHttpUrl = "v1/organizations/accounts/{account_id}"
		getAccountProduct = "organizations"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{account_id}", state.Primary.ID)

	getAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountResp, err := getAccountClient.Request("GET", getAccountPath, &getAccountOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Account: %s", err)
	}

	getAccountRespBody, err := utils.FlattenResponse(getAccountResp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("account.status", getAccountRespBody, "").(string)
	if status == "" || status == "pending_closure" || status == "suspended" {
		return nil, golangsdk.ErrDefault404{}
	}
	return getAccountRespBody, nil
}

func TestAccAccount_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_organizations_account.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getAccountResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsOrganizationalUnitId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAccount_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "parent_id", acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID),
					resource.TestCheckResourceAttr(rName, "phone", "13245678978"),
					resource.TestCheckResourceAttr(rName, "agency_name", "OrganizationAccountAccessAgency"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
				),
			},
			{
				Config: testAccAccount_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"data.huaweicloud_organizations_organization.test", "root_id"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform1"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"agency_name"},
			},
		},
	})
}

func testAccAccount_basic_step1(name string) string {
	return fmt.Sprintf(`
# 'email' parameter is available when creating account in international website.
# Create an account in the specified organizational unit.
resource "huaweicloud_organizations_account" "test" {
  name        = "%[1]s"
  parent_id   = "%[2]s"
  phone       = "13245678978"
  agency_name = "OrganizationAccountAccessAgency"
  description = "Created by terraform script"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, name, acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID)
}

func testAccAccount_basic_step2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_account" "test" {
  name = "%[1]s"
  # Move the account to the root organizational unit.
  parent_id   = data.huaweicloud_organizations_organization.test.root_id
  phone       = "13245678978"
  agency_name = "OrganizationAccountAccessAgency"

  tags = {
    foo1  = "bar"
    owner = "terraform1"
  }
}
`, name)
}
