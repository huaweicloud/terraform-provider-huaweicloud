package identitycenter

import (
	"fmt"
	"regexp"
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
					resource.TestCheckResourceAttr(rName, "phone_number", "13600000000"),
					resource.TestCheckResourceAttr(rName, "user_type", "test_user_type"),
					resource.TestCheckResourceAttr(rName, "title", "test_title"),
					resource.TestCheckResourceAttr(rName, "addresses.#", "1"),
					resource.TestCheckResourceAttr(rName, "addresses.0.country", "test_country"),
					resource.TestCheckResourceAttr(rName, "addresses.0.formatted", "test_formatted"),
					resource.TestCheckResourceAttr(rName, "addresses.0.locality", "test_locality"),
					resource.TestCheckResourceAttr(rName, "addresses.0.postal_code", "test_postal_code"),
					resource.TestCheckResourceAttr(rName, "addresses.0.region", "test_region"),
					resource.TestCheckResourceAttr(rName, "addresses.0.street_address", "test_street_address"),
					resource.TestCheckResourceAttr(rName, "enterprise.#", "1"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.cost_center", "test_cost_center"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.department", "test_department"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.division", "test_division"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.employee_number", "test_employee_number"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.organization", "test_organization"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.manager", "test_manager"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "email_verified"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
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
					resource.TestCheckResourceAttr(rName, "phone_number", "13600000001"),
					resource.TestCheckResourceAttr(rName, "user_type", "test_user_type_update"),
					resource.TestCheckResourceAttr(rName, "title", "test_title_update"),
					resource.TestCheckResourceAttr(rName, "addresses.#", "1"),
					resource.TestCheckResourceAttr(rName, "addresses.0.country", "test_country_update"),
					resource.TestCheckResourceAttr(rName, "addresses.0.formatted", "test_formatted_update"),
					resource.TestCheckResourceAttr(rName, "addresses.0.locality", "test_locality_update"),
					resource.TestCheckResourceAttr(rName, "addresses.0.postal_code", "test_postal_code_update"),
					resource.TestCheckResourceAttr(rName, "addresses.0.region", "test_region_update"),
					resource.TestCheckResourceAttr(rName, "addresses.0.street_address", "test_street_address_update"),
					resource.TestCheckResourceAttr(rName, "enterprise.#", "1"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.cost_center", "test_cost_center_update"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.department", "test_department_update"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.division", "test_division_update"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.employee_number", "test_employee_number_update"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.organization", "test_organization_update"),
					resource.TestCheckResourceAttr(rName, "enterprise.0.manager", "test_manager_update"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "email_verified"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityCenterUserImportState(rName),
				ImportStateVerifyIgnore: []string{"password_mode", "user_status"},
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
  phone_number      = "13600000000"
  user_type         = "test_user_type"
  title             = "test_title"
  enabled           = "false"

  addresses {
    country        = "test_country"
    formatted      = "test_formatted"
    locality       = "test_locality"
    postal_code    = "test_postal_code"
    region         = "test_region"
    street_address = "test_street_address"
  }

  enterprise {
    cost_center     = "test_cost_center"
    department      = "test_department"
    division        = "test_division"
    employee_number = "test_employee_number"
    organization    = "test_organization"
    manager         = "test_manager"
  }

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
  phone_number      = "13600000001"
  user_type         = "test_user_type_update"
  title             = "test_title_update"
  enabled           = "true"

  addresses {
    country        = "test_country_update"
    formatted      = "test_formatted_update"
    locality       = "test_locality_update"
    postal_code    = "test_postal_code_update"
    region         = "test_region_update"
    street_address = "test_street_address_update"
  }

  enterprise {
    cost_center     = "test_cost_center_update"
    department      = "test_department_update"
    division        = "test_division_update"
    employee_number = "test_employee_number_update"
    organization    = "test_organization_update"
    manager         = "test_manager_update"
  }
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
