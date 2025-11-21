package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataExtendFlavors_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dms_kafka_extend_flavors.all"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byType   = "data.huaweicloud_dms_kafka_extend_flavors.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byChargingMode   = "data.huaweicloud_dms_kafka_extend_flavors.filter_by_charging_mode"
		dcByChargingMode = acceptance.InitDataSourceCheck(byChargingMode)

		byArchType   = "data.huaweicloud_dms_kafka_extend_flavors.filter_by_arch_type"
		dcByArchType = acceptance.InitDataSourceCheck(byArchType)

		byStorageSpecCode   = "data.huaweicloud_dms_kafka_extend_flavors.filter_by_storage_spec_code"
		dcByStorageSpecCode = acceptance.InitDataSourceCheck(byStorageSpecCode)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataExtendFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "versions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSource, "flavors.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_useful", "true"),
					dcByChargingMode.CheckResourceExists(),
					resource.TestCheckOutput("is_charging_mode_useful", "true"),
					dcByArchType.CheckResourceExists(),
					resource.TestCheckOutput("is_arch_type_useful", "true"),
					dcByStorageSpecCode.CheckResourceExists(),
					resource.TestCheckOutput("is_storage_spec_code_useful", "true"),
				),
			},
		},
	})
}

func testAccDataExtendFlavors_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instances" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_dms_kafka_extend_flavors" "all" {
  instance_id = "%[1]s"
}

# Filter by instance type.
locals {
  instance_type = try(data.huaweicloud_dms_kafka_instances.test.instances[0].type, null)
}

data "huaweicloud_dms_kafka_extend_flavors" "filter_by_type" {
  instance_id = "%[1]s"
  type        = local.instance_type
}

locals {
  filter_by_type_result = [for v in data.huaweicloud_dms_kafka_extend_flavors.filter_by_type.flavors[*].type : v == local.instance_type]
}

output "is_type_useful" {
  value = length(local.filter_by_type_result) > 0 && alltrue(local.filter_by_type_result)
}

# Filter by instance charging mode.
locals {
  charging_mode = try(data.huaweicloud_dms_kafka_instances.test.instances[0].charging_mode, null)
}

data "huaweicloud_dms_kafka_extend_flavors" "filter_by_charging_mode" {
  instance_id   = "%[1]s"
  charging_mode = local.charging_mode
}

locals {
  filter_by_charging_mode_result = [for v in data.huaweicloud_dms_kafka_extend_flavors.filter_by_charging_mode.flavors[*].charging_modes :
  contains(v, local.charging_mode)]
}

output "is_charging_mode_useful" {
  value = length(local.filter_by_charging_mode_result) > 0 && alltrue(local.filter_by_charging_mode_result)
}

# Filter by instance architecture type.
locals {
  arch_type = try(data.huaweicloud_dms_kafka_extend_flavors.all.flavors[0].arch_types[0], null)
}

data "huaweicloud_dms_kafka_extend_flavors" "filter_by_arch_type" {
  instance_id = "%[1]s"
  arch_type   = local.arch_type
}

locals {
  filter_by_arch_type_result = [for v in data.huaweicloud_dms_kafka_extend_flavors.filter_by_arch_type.flavors[*].arch_types :
  contains(v, local.arch_type)]
}

output "is_arch_type_useful" {
  value = length(local.filter_by_arch_type_result) > 0 && alltrue(local.filter_by_arch_type_result)
}

# Filter by instance storage spec code.
locals {
  storage_spec_code = try(data.huaweicloud_dms_kafka_instances.test.instances[0].storage_spec_code, null)
}

data "huaweicloud_dms_kafka_extend_flavors" "filter_by_storage_spec_code" {
  instance_id       = "%[1]s"
  storage_spec_code = local.storage_spec_code
}

locals {
  filter_by_storage_spec_code_result = [
    for v in data.huaweicloud_dms_kafka_extend_flavors.filter_by_storage_spec_code.flavors[*].ios[*].storage_spec_code :
    contains(v, local.storage_spec_code)
  ]
}

output "is_storage_spec_code_useful" {
  value = length(local.filter_by_storage_spec_code_result) > 0 && alltrue(local.filter_by_storage_spec_code_result)
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
