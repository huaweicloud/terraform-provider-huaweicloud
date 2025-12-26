package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5UserAttachedPolicies_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identityv5_user_attached_policies.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5UserAttachedPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "attached_policies.#"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5UserAttachedPolicies_basic() string {
	return `
data "huaweicloud_identityv5_user_attached_policies" "test" {
  user_id = "bdbd75fde59f49eea1b3ea1d2426f4d9"
}
`
}
