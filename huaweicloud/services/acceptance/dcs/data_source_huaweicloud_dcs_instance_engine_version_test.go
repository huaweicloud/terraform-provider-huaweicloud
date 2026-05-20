package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsInstanceEngineVersion_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_instance_engine_version.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsInstanceEngineVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "engine_minor_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "latest_engine_minor_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "engine_minor_version_upgradable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxy_minor_version_upgradable"),
				),
			},
		},
	})
}

func testAccDataSourceDcsInstanceEngineVersion_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_instance_engine_version" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
}
`, testAccDcsInstance_base(name))
}
