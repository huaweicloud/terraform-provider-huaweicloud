package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationalUnits_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_organizational_units.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationalUnits_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_organizational_units.#"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_organizational_units.0.manage_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_organizational_units.0.organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_organizational_units.0.organizational_unit_type"),
				),
			},
		},
	})
}

const testAccDataSourceOrganizationalUnits_basic = `
data "huaweicloud_rgc_organizational_units" "test" {}
`
