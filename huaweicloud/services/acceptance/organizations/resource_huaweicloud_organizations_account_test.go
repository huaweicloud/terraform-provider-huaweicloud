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
		name      = acceptance.RandomAccResourceName()
		withPhone = acceptance.RandomAccResourceName()

		obj interface{}

		rName = "huaweicloud_organizations_account.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getAccountResourceFunc)

		rNameWithPhone = "huaweicloud_organizations_account.with_phone"
		rcWithPhone    = acceptance.InitResourceCheck(rNameWithPhone, &obj, getAccountResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithPhone.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAccount_basic_step1(name, withPhone),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_organizations_organizational_unit.test", "id"),
					resource.TestCheckResourceAttr(rName, "agency_name", "OrganizationAccountAccessAgency"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
					rcWithPhone.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithPhone, "name", withPhone),
					resource.TestCheckResourceAttrPair(rNameWithPhone, "parent_id",
						"huaweicloud_organizations_organizational_unit.test", "id"),
					resource.TestCheckResourceAttr(rNameWithPhone, "phone", "13245678978"),
					resource.TestCheckResourceAttr(rNameWithPhone, "agency_name", "OrganizationAccountAccessAgency"),
					resource.TestCheckResourceAttr(rNameWithPhone, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rNameWithPhone, "tags.%", "2"),
					resource.TestCheckResourceAttr(rNameWithPhone, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rNameWithPhone, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(rNameWithPhone, "urn"),
					resource.TestCheckResourceAttrSet(rNameWithPhone, "joined_at"),
					resource.TestCheckResourceAttrSet(rNameWithPhone, "joined_method"),
				),
			},
			{
				Config: testAccAccount_basic_step2(name, withPhone),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"data.huaweicloud_organizations_organization.test", "root_id"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform1"),
					rcWithPhone.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithPhone, "name", withPhone),
					resource.TestCheckResourceAttrPair(rNameWithPhone, "parent_id",
						"data.huaweicloud_organizations_organization.test", "root_id"),
					resource.TestCheckResourceAttr(rNameWithPhone, "phone", "13245678978"),
					resource.TestCheckResourceAttr(rNameWithPhone, "agency_name", "OrganizationAccountAccessAgency"),
					resource.TestCheckResourceAttr(rNameWithPhone, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(rNameWithPhone, "tags.%", "0"),
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

func testAccAccount_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[1]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}
`, name)
}

func testAccAccount_basic_step1(name, withPhone string) string {
	return fmt.Sprintf(`
%[1]s

# Create two accounts in the specified organizational unit.
resource "huaweicloud_organizations_account" "test" {
  name        = "%[2]s"
  parent_id   = huaweicloud_organizations_organizational_unit.test.id
  agency_name = "OrganizationAccountAccessAgency"
  description = "Created by terraform script"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}

# 'email' parameter is available and required when creating account in international website.
resource "huaweicloud_organizations_account" "with_phone" {
  name        = "%[3]s"
  parent_id   = huaweicloud_organizations_organizational_unit.test.id
  phone       = "13245678978"
  agency_name = "OrganizationAccountAccessAgency"
  description = "Created by terraform script"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, testAccAccount_basic_base(name), name, withPhone)
}

func testAccAccount_basic_step2(name, withPhone string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_account" "test" {
  name        = "%[2]s"
  # Move the account to the root organizational unit.
  parent_id   = data.huaweicloud_organizations_organization.test.root_id
  agency_name = "OrganizationAccountAccessAgency"

  tags = {
    foo1  = "bar"
    owner = "terraform1"
  }
}

resource "huaweicloud_organizations_account" "with_phone" {
  name        = "%[3]s"
  # Move the account to the root organizational unit.
  parent_id   = data.huaweicloud_organizations_organization.test.root_id
  phone       = "13245678978"
  agency_name = "OrganizationAccountAccessAgency"
  description = "Updated by terraform script"
}
`, testAccAccount_basic_base(name), name, withPhone)
}
