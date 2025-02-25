package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitmqExtendFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rabbitmq_extend_flavors.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsRabbitmqExtendFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckOutput("type_validation", "true"),
					resource.TestCheckOutput("arch_types_validation", "true"),
					resource.TestCheckOutput("charging_modes_validation", "true"),
					resource.TestCheckOutput("storage_spec_code_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsRabbitmqExtendFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_rabbitmq_extend_flavors" "all" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
}

data "huaweicloud_dms_rabbitmq_extend_flavors" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_instance.test]

  instance_id       = huaweicloud_dms_rabbitmq_instance.test.id
  type              = local.test_refer.type
  arch_type         = local.test_refer.arch_types[0]
  charging_mode     = local.test_refer.charging_modes[0]
  storage_spec_code = local.test_refer.ios[0].storage_spec_code
}

locals {
  test_refer   = data.huaweicloud_dms_rabbitmq_extend_flavors.all.flavors[0]
  test_results = data.huaweicloud_dms_rabbitmq_extend_flavors.test
}

output "type_validation" {
  value = alltrue([for a in local.test_results.flavors[*].type : a == local.test_refer.type])
}

output "arch_types_validation" {
  value = alltrue([for a in local.test_results.flavors[*].arch_types : contains(a, local.test_refer.arch_types[0])])
}

output "charging_modes_validation" {
  value = alltrue([for a in local.test_results.flavors[*].charging_modes : contains(a, local.test_refer.charging_modes[0])])
}

output "storage_spec_code_validation" {
  value = alltrue([for ios in local.test_results.flavors[*].ios :
    alltrue([for io in ios : io.storage_spec_code == local.test_refer.ios[0].storage_spec_code])])
}
`, testAccDmsRabbitmqInstance_newFormat_single(name))
}
