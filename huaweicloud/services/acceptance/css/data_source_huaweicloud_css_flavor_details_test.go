package css

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFlavorDetails_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_flavor_details.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFlavorDetails_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "str_id"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "cond_operation_status"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_detail.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_detail.0.value"),
				),
			},
		},
	})
}

const testAccDataSourceFlavorDetails_basic = `
data "huaweicloud_css_flavors" "test" {}

data "huaweicloud_css_flavor_details" "test" {
  flavor_id = data.huaweicloud_css_flavors.test.flavors[0].id
}
`
