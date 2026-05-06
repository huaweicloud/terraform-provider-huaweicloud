package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipBindings_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_bindings.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalEipBindings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "geip_bindings.#"),
				),
			},
		},
	})
}

func testAccDataSourceGlobalEipBindings_basic() string {
	return `
data "huaweicloud_global_eip_bindings" "test" {
  fields   = ["geip_id", "geip_ip_address"]
  sort_key = "geip_id"
  sort_dir = "asc"
}
`
}
