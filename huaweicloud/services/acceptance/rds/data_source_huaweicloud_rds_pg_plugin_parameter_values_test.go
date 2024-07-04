package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPgPluginParameterValues_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_pg_plugin_parameter_values.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgPluginParameterValues_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "values.0",
						"data.huaweicloud_rds_pg_plugin_parameter_value_range.test", "values.0"),
				),
			},
		},
	})
}

func testAccDatasourcePgPluginParameterValues_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_pg_plugin_parameter_values" "test" {
  depends_on = [huaweicloud_rds_pg_plugin_parameter.test]

  instance_id = huaweicloud_rds_instance.test.id
  name        = "shared_preload_libraries"
}
`, testPgPluginParameter_basic(name))
}
