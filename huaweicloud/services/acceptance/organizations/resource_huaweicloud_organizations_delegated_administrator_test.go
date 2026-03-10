package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/organizations"
)

func getDelegatedAdministratorResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("organizations", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	return organizations.GetDelegatedAdministrator(client, state.Primary.Attributes["account_id"],
		state.Primary.Attributes["service_principal"])
}

func TestAccDelegatedAdministrator_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_organizations_delegated_administrator.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDelegatedAdministratorResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsAccountName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDelegatedAdministrator_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "account_id",
						"data.huaweicloud_organizations_accounts.test", "accounts.0.id"),
					resource.TestCheckResourceAttrPair(rName, "service_principal",
						"huaweicloud_organizations_trusted_service.test", "service"),
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

func testDelegatedAdministrator_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_trusted_service" "test" {
  service = "service.SecMaster"
}

data "huaweicloud_organizations_accounts" "test" {
  name = "%[1]s"
}

resource "huaweicloud_organizations_delegated_administrator" "test" {
  account_id        = try(data.huaweicloud_organizations_accounts.test.accounts[0].id, "")
  service_principal = huaweicloud_organizations_trusted_service.test.service
}
`, acceptance.HW_ORGANIZATIONS_ACCOUNT_NAME)
}
