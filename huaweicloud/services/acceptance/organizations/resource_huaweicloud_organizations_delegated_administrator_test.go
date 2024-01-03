package organizations

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDelegatedAdministratorResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getDelegatedAdministrator: Query Organizations delegated administrator
	var (
		getDelegatedAdministratorHttpUrl = "v1/organizations/delegated-administrators"
		getDelegatedAdministratorProduct = "organizations"
	)
	getDelegatedAdministratorClient, err := cfg.NewServiceClient(getDelegatedAdministratorProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Organizations client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <account_id>/<service_principal>")
	}
	accountID := parts[0]
	servicePrincipal := parts[1]

	getDelegatedAdministratorPath := getDelegatedAdministratorClient.Endpoint + getDelegatedAdministratorHttpUrl
	getDelegatedAdministratorQueryParams := buildGetDelegatedAdministratorQueryParams(servicePrincipal)
	getDelegatedAdministratorPath += getDelegatedAdministratorQueryParams

	getDelegatedAdministratorResp, err := pagination.ListAllItems(
		getDelegatedAdministratorClient,
		"marker",
		getDelegatedAdministratorPath,
		&pagination.QueryOpts{MarkerField: "account_id"})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations delegated administrator: %s", err)
	}

	getDelegatedAdministratorRespJson, err := json.Marshal(getDelegatedAdministratorResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations delegated administrator: %s", err)
	}
	var getDelegatedAdministratorRespBody interface{}
	err = json.Unmarshal(getDelegatedAdministratorRespJson, &getDelegatedAdministratorRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Organizations delegated administrator: %s", err)
	}

	delegatedAdministrator := utils.PathSearch(fmt.Sprintf("delegated_administrators|[?account_id=='%s']|[0]",
		accountID), getDelegatedAdministratorRespBody, nil)
	if delegatedAdministrator == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return delegatedAdministrator, nil
}

func buildGetDelegatedAdministratorQueryParams(servicePrincipal string) string {
	return fmt.Sprintf("?service_principal=%v", servicePrincipal)
}

func TestAccDelegatedAdministrator_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_organizations_delegated_administrator.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDelegatedAdministratorResourceFunc,
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
				Config: testDelegatedAdministrator_basic(acceptance.HW_ORGANIZATIONS_ACCOUNT_NAME),
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

func testDelegatedAdministrator_basic(name string) string {
	return fmt.Sprintf(`

%[1]s

data "huaweicloud_organizations_accounts" "test" {
  name = "%[2]s"
}

resource "huaweicloud_organizations_delegated_administrator" "test" {
  account_id        = data.huaweicloud_organizations_accounts.test.accounts.0.id
  service_principal = huaweicloud_organizations_trusted_service.test.service
}
`, testTrustedService_basic(), name)
}
