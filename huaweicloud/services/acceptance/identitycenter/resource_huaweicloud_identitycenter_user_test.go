package identitycenter

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

func getIdentityCenterUserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getIdentityCenterUser: Query Identity Center user
	var (
		getIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}"
		getIdentityCenterUserProduct = "identitystore"
	)
	getIdentityCenterUserClient, err := cfg.NewServiceClient(getIdentityCenterUserProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	getIdentityCenterUserPath := getIdentityCenterUserClient.Endpoint + getIdentityCenterUserHttpUrl
	getIdentityCenterUserPath = strings.ReplaceAll(getIdentityCenterUserPath, "{identity_store_id}",
		state.Primary.Attributes["identity_store_id"])
	getIdentityCenterUserPath = strings.ReplaceAll(getIdentityCenterUserPath, "{user_id}", state.Primary.ID)

	getIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterUserResp, err := getIdentityCenterUserClient.Request("GET", getIdentityCenterUserPath,
		&getIdentityCenterUserOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center User: %s", err)
	}
	return utils.FlattenResponse(getIdentityCenterUserResp)
}

func TestAccIdentityCenterUser_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_identitycenter_user.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterUserResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterUser_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "user_name", name),
					resource.TestCheckResourceAttr(rName, "family_name", "test_family_name"),
					resource.TestCheckResourceAttr(rName, "given_name", "test_given_name"),
					resource.TestCheckResourceAttr(rName, "display_name", "test_display_name"),
					resource.TestCheckResourceAttr(rName, "email", "email@example.com"),
				),
			},
			{
				Config: testIdentityCenterUser_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "user_name", name),
					resource.TestCheckResourceAttr(rName, "family_name", "test_family_name_update"),
					resource.TestCheckResourceAttr(rName, "given_name", "test_given_name_update"),
					resource.TestCheckResourceAttr(rName, "display_name", "test_display_name_update"),
					resource.TestCheckResourceAttr(rName, "email", "email_update@example.com"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityCenterUserImportState(rName),
				ImportStateVerifyIgnore: []string{"password_mode"},
			},
		},
	})
}

func testIdentityCenterUser_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_user" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  user_name         = "%s"
  password_mode     = "OTP"
  family_name       = "test_family_name"
  given_name        = "test_given_name"
  display_name      = "test_display_name"
  email             = "email@example.com"
}
`, testAccDatasourceIdentityCenter_basic(), name)
}

func testIdentityCenterUser_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_user" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  user_name         = "%s"
  password_mode     = "OTP"
  family_name       = "test_family_name_update"
  given_name        = "test_given_name_update"
  display_name      = "test_display_name_update"
  email             = "email_update@example.com"
}
`, testAccDatasourceIdentityCenter_basic(), name)
}

// testIdentityCenterUserImportState use to return an id with format <identity_store_id>/<user_id>
func testIdentityCenterUserImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		return identityStoreId + "/" + rs.Primary.ID, nil
	}
}
