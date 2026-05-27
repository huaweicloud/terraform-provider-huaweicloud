package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBPtApplicableInstances_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.huaweicloud_geminidb_pt_applicable_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBPtApplicableInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),

					resource.TestCheckOutput("instance_id_filter_useful", "true"),
					resource.TestCheckOutput("instance_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBPtApplicableInstances_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_geminidb_parameter_template" "test" {
  name        = "%[2]s"
  description = "configuration based on instance"
  instance_id = huaweicloud_geminidb_instance.test.id

  depends_on = [huaweicloud_geminidb_instance.test]
}

data "huaweicloud_geminidb_pt_applicable_instances" "test" {
  config_id = huaweicloud_geminidb_parameter_template.test.id

  depends_on = [huaweicloud_geminidb_parameter_template.test]
}

data "huaweicloud_geminidb_pt_applicable_instances" "instance_id_filter" {
  config_id   = huaweicloud_geminidb_parameter_template.test.id
  instance_id = huaweicloud_geminidb_instance.test.id

  depends_on = [huaweicloud_geminidb_parameter_template.test]
}

output "instance_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_pt_applicable_instances.instance_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_pt_applicable_instances.instance_id_filter.instances[*].id : v == huaweicloud_geminidb_instance.test.id]
  )
}

data "huaweicloud_geminidb_pt_applicable_instances" "instance_name_filter" {
  config_id     = huaweicloud_geminidb_parameter_template.test.id
  instance_name = huaweicloud_geminidb_instance.test.name

  depends_on = [huaweicloud_geminidb_parameter_template.test]
}

output "instance_name_filter_useful" {
  value = length(data.huaweicloud_geminidb_pt_applicable_instances.instance_name_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_pt_applicable_instances.instance_name_filter.instances[*].name : v == huaweicloud_geminidb_instance.test.name]
  )
}
`, testAccGeminiDbInstance_basic(rName), rName)
}
