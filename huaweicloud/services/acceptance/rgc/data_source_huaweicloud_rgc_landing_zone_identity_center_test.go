package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLandingZoneIdentityCenter_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_landing_zone_identity_center.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLandingZoneIdentityCenterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "identity_store_id"),
					resource.TestCheckResourceAttrSet(dataSource, "user_portal_url"),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.0.permission_set_id"),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.0.permission_set_name"),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "groups.0.description"),
				),
			},
		},
	})
}

const testAccDataSourceLandingZoneIdentityCenterConfig_basic = `
data "huaweicloud_rgc_landing_zone_identity_center" "test" {}
`
