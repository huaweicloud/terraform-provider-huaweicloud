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

func getAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	// getAccount: Query RGC account via organizations API
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
	rName := "huaweicloud_rgc_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOrganization(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccount_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"email", "phone", "identity_store_user_name", "identity_store_email",
					"parent_organizational_unit_name", "parent_organizational_unit_id"},
			},
		},
	})
}

func testAccount_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_account" "test" {
  name                            = "%[1]s"
  email                           = "%[1]s@terraform.com"
  phone                           = "13987654321"
  identity_store_user_name        = "%[1]s"
  identity_store_email            = "%[1]s@terraform.com"
  parent_organizational_unit_name = "%[2]s"
  parent_organizational_unit_id   = "%[3]s"
}
`, name, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_NAME, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}

func TestAccAccount_blueprint(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rgc_account.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccountResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOrganization(t)
			acceptance.TestAccPreCheckRGCBlueprint(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccount_blueprint(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"email", "phone", "identity_store_user_name", "identity_store_email",
					"parent_organizational_unit_name", "parent_organizational_unit_id", "blueprint"},
			},
		},
	})
}

func testAccount_blueprint(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rgc_account" "test" {
  name                            = "%[1]s"
  email                           = "%[1]s@terraform.com"
  phone                           = "13987654321"
  identity_store_user_name        = "%[1]s"
  identity_store_email            = "%[1]s@terraform.com"
  parent_organizational_unit_name = "%[2]s"
  parent_organizational_unit_id   = "%[3]s"

  blueprint {
    blueprint_product_id                    = "%[4]s"
    blueprint_product_version               = "%[5]s"
    is_blueprint_has_multi_account_resource = false
  }
}
`, name, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_NAME, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID,
		acceptance.HW_RGC_BLUEPRINT_PRODUCT_ID, acceptance.HW_RGC_BLUEPRINT_PRODUCT_VERSION)
}
