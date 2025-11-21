package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5Groups_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_groups.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5Groups_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.#"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5Groups_basic() string {
	return `
data "huaweicloud_identityv5_groups" "test" {}
`
}
