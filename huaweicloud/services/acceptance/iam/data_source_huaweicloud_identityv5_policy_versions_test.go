package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityV5PolicyVersions_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identityv5_policy_versions.test"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIdentityV5PolicyVersions_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.#"),
				),
			},
			{
				Config: testAccDataSourceIdentityV5PolicyVersionsWithVersionId_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "versions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "versions.0.version_id", "v1"),
					resource.TestCheckResourceAttr(dataSourceName, "versions.0.is_default", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "versions.0.document"),
				),
			},
		},
	})
}

var testAccDataSourceIdentityV5PolicyVersions_basic = `
data "huaweicloud_identityv5_policy_versions" "test" {
  policy_id = "NATReadOnlyPolicy"
}
`

var testAccDataSourceIdentityV5PolicyVersionsWithVersionId_basic = `
data "huaweicloud_identityv5_policy_versions" "test" {
  policy_id  = "NATReadOnlyPolicy"
  version_id = "v1"
}
`
