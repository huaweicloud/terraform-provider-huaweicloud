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
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_organizations_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountResourceFunc,
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
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test_account_description"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
				),
			},
			{
				Config: testAccount_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "parent_id",
						acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
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

func testAccount_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account" "test" {
  name        = "%s"
  email       = "account_1@abc.com"
  phone       = "13987654321"
  description = "test_account_description"

  tags = {
    "key1" = "value1"
    "key2" = "value2"
  }
}
`, name)
}

func testAccount_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_account" "test" {
  name        = "%s"
  email       = "account_1@abc.com"
  phone       = "13987654321"
  description = ""
  parent_id   = "%s"

  tags = {
    "key3" = "value3"
    "key4" = "value4"
  }
}
`, name, acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID)
}
