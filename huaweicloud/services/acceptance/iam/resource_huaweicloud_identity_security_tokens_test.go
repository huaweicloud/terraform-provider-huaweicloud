package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"testing"
)

func TestAccIdentitySecurityTokens_basic(t *testing.T) {
	resourceName := "huaweicloud_identity_security_tokens.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AccIdentitySecurityTokens_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "credential.#"),
				),
			},
		},
	})
}

const AccIdentitySecurityTokens_basic = `
resource "huaweicloud_identity_security_tokens" "test" {
  version = "1.1"
  action = ["vpc:ports:create"]
  effect = "Allow"
}
`
