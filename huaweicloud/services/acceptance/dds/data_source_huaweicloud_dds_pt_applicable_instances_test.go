package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsPtApplicableInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_pt_applicable_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDdsPtApplicableInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.entities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.entities.0.entity_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.entities.0.entity_name"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDdsPtApplicableInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_dds_pt_applicable_instances" "test" {
  depends_on = [huaweicloud_dds_instance.instance]

  configuration_id = huaweicloud_dds_parameter_template.test.id
}
`, testAccDDSInstanceV3Config_basic(name, 8800), testDdsParameterTemplate_basic(name))
}
