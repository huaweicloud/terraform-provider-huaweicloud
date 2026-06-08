package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstalledPlugins_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_instance_installed_plugins.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(
		t, resource.TestCase{
			PreCheck: func() {
				acceptance.TestAccPreCheck(t)
				acceptance.TestAccPreCheckGaussDBInstanceId(t)
			},
			ProviderFactories: acceptance.TestAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testDataSourceInstalledPlugins_basic(),
					Check: resource.ComposeTestCheckFunc(
						dc.CheckResourceExists(),
						resource.TestCheckResourceAttrSet(dataSource, "plugins.#"),
					),
				},
			},
		},
	)
}

func testDataSourceInstalledPlugins_basic() string {
	return fmt.Sprintf(`

data "huaweicloud_gaussdb_instance_installed_plugins" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
