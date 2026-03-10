package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatastoreFlavors_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_datastore_flavors.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSDatastoreId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatastoreFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datastore_id_str"),
					resource.TestCheckResourceAttrSet(dataSource, "dbname"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.typename"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.diskrange"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.flavors.0.cond_operation_status"),
					resource.TestCheckResourceAttrSet(dataSource, "model_list.#"),

					resource.TestCheckOutput("datastore_version_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDatastoreFlavors_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_datastore_flavors" "test" {
  datastore_id = "%[1]s"
}

locals {
  datastore_version_id = data.huaweicloud_css_datastore_flavors.test.versions[0].id
}

data "huaweicloud_css_datastore_flavors" "datastore_version_id_filter" {
  datastore_id         = "%[1]s"	
  datastore_version_id = local.datastore_version_id
}

output "datastore_version_id_filter_useful" {
  value = length(data.huaweicloud_css_datastore_flavors.datastore_version_id_filter.versions) > 0
}
`, acceptance.HW_CSS_DATASTORE_ID)
}
