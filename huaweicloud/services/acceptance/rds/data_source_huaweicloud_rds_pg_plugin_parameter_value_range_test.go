package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePgPluginParameterValueRange_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_pg_plugin_parameter_value_range.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgPluginParameterValueRange_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "restart_required"),
					resource.TestCheckResourceAttrSet(rName, "values.#"),
					resource.TestCheckResourceAttrSet(rName, "default_values.#"),
				),
			},
		},
	})
}

func testAccDatasourcePgPluginParameterValueRange_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_pg_plugin_parameter_value_range" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "shared_preload_libraries"
}
`, testAccRdsInstance_basic(name))
}
