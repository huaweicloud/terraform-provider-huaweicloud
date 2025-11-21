package iam_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5AuthorizationSchema_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_authorization_schema.schema"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5AuthorizationSchema_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet("data.huaweicloud_identityv5_authorization_schema.schema", "actions.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_identityv5_authorization_schema.schema", "conditions.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_identityv5_authorization_schema.schema", "operations.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_identityv5_authorization_schema.schema", "resources.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_identityv5_authorization_schema.schema", "version"),
				),
			},
		},
	})
}

func testAccDataSourceIdentityV5AuthorizationSchema_basic() string {
	return `
data "huaweicloud_identityv5_authorization_schema" "schema" {
	service_code = "iam"	
}
`
}
