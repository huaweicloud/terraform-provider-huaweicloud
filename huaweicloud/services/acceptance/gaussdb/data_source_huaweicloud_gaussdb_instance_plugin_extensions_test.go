package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstancePluginExtensions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_instance_plugin_extensions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstancePluginExtensions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "extensions.#"),
				),
			},
		},
	},
	)
}

func testDataSourceInstancePluginExtensions_basic() string {
	return fmt.Sprintf(`

data "huaweicloud_gaussdb_instance_plugin_extensions" "test" {
  instance_id = "%[1]s"
  db_name     = "gauss-6ea"
  plugin_name = "postgis"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
