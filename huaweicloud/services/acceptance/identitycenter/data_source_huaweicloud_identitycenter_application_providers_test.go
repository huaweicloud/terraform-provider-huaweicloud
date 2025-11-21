package identitycenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterApplicationProviders_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_application_providers.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterApplicationProviders_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName,
						"application_providers.0.application_provider_urn", "IdentityCenter:::applicationProvider:custom-saml"),
					resource.TestCheckResourceAttr(rName,
						"application_providers.0.display_data.0.description", "Custom SAML 2.0 application"),
					resource.TestCheckResourceAttr(rName,
						"application_providers.0.display_data.0.display_name", "Custom SAML 2.0 application"),
					resource.TestCheckResourceAttr(rName, "application_providers.0.display_data.0.icon_url", ""),
					resource.TestCheckResourceAttr(rName, "application_providers.0.federation_protocol", "SAML")),
			},
		},
	})
}

const testAccDatasourceIdentityCenterApplicationProviders_basic = `
data "huaweicloud_identitycenter_application_providers" "test" {}
`
