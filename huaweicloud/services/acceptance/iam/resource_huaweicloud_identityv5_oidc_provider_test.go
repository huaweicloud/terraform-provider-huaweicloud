package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5OidcProviderFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5OidcProviderById(client, state.Primary.ID)
}

// Please ensure that the user executing the acceptance test has 'admin' permission,
// and the account has enabled the IAM v5 (new IAM) feature.
func TestAccV5OidcProvider_basic(t *testing.T) {
	var (
		obj    interface{}
		rName  = "huaweicloud_identityv5_oidc_provider.test"
		dsName = "data.huaweicloud_identityv5_oidc_providers.test"
		rc     = acceptance.InitResourceCheck(rName, &obj, getV5OidcProviderFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckIAMV5(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5OidcProvider_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "acc-oidc-provider"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "client_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "client_ids.0", "sts.amazonaws.com"),
					resource.TestCheckResourceAttrSet(rName, "url"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
				),
			},
			{
				Config: testAccV5OidcProvider_basic_step2() + testAccV5OidcProviders_dataSource(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(rName, "client_ids.#", "2"),
					resource.TestCheckResourceAttr(rName, "thumbprints.#", "1"),
					resource.TestCheckResourceAttr(dsName, "providers.#", "1"),
					resource.TestCheckResourceAttr(dsName, "providers.0.name", "acc-oidc-provider"),
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

func testAccV5OidcProvider_basic_step1() string {
	return `
resource "huaweicloud_identityv5_oidc_provider" "test" {
  name        = "acc-oidc-provider"
  url         = "https://token.actions.githubusercontent.com"
  description = "Created by terraform script"
  client_ids  = ["sts.amazonaws.com"]
}
`
}

func testAccV5OidcProvider_basic_step2() string {
	return `
resource "huaweicloud_identityv5_oidc_provider" "test" {
  name        = "acc-oidc-provider"
  url         = "https://token.actions.githubusercontent.com"
  description = "Updated by terraform script"
  client_ids  = ["sts.amazonaws.com", "another-client-id"]
  thumbprints = ["b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5b5"]
}
`
}

func testAccV5OidcProviders_dataSource() string {
	return `
data "huaweicloud_identityv5_oidc_providers" "test" {
  name       = "acc-oidc-provider"
  depends_on = [huaweicloud_identityv5_oidc_provider.test]
}
`
}
