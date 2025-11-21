package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5ServicePrincipals_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5ServicePrincipals_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.huaweicloud_identityv5_service_principals.test", "service_principals.#"),
				),
			},
		},
	})
}

var testAccDataSourceIdentityV5ServicePrincipals_basic = `
data "huaweicloud_identityv5_service_principals" "test" {}
`
