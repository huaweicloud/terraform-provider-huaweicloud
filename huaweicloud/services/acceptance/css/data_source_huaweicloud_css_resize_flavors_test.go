package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResizeFlavors_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_resize_flavors.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResizeFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_id"),
					resource.TestCheckResourceAttrSet(dataSource, "dbname"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.typename"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.diskrange"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.cond_operation_status"),

					resource.TestCheckOutput("type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceResizeFlavors_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_resize_flavors" "test" {
  cluster_id = "%[1]s"
}

locals {
  type = data.huaweicloud_css_resize_flavors.test.versions[0].flavors[0].typename
}

data "huaweicloud_css_resize_flavors" "type_filter" {
  cluster_id = "%[1]s"	
  type       = local.type
}

output "type_filter_useful" {
  value = length(data.huaweicloud_css_resize_flavors.type_filter.versions) > 0 && alltrue(
    [for v in data.huaweicloud_css_resize_flavors.type_filter.versions[0].flavors[*].typename : v == local.type]
  )
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
