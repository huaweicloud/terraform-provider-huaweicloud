package identitycenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceServiceProviderConfiguration_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_service_provider_configuration.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceServiceProviderConfiguration_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "sp_oidc_config.#"),
					resource.TestCheckResourceAttrSet(rName, "sp_oidc_config.0.redirect_url"),
					resource.TestCheckResourceAttrSet(rName, "sp_saml_config.#"),
					resource.TestCheckResourceAttrSet(rName, "sp_saml_config.0.acs_url"),
					resource.TestCheckResourceAttrSet(rName, "sp_saml_config.0.issuer"),
					resource.TestCheckResourceAttrSet(rName, "sp_saml_config.0.metadata"),
				),
			},
		},
	})
}

const testAccDatasourceServiceProviderConfiguration_basic = `
data "huaweicloud_identitycenter_instance" "test" {}

data "huaweicloud_identitycenter_service_provider_configuration" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`
