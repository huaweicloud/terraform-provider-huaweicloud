package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSocPreprocessRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_secmaster_soc_preprocess_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSocPreprocessRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.mapping_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.mapper_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.mapper_type_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data.0.action"),

					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("mapping_id_filter_useful", "true"),
					resource.TestCheckOutput("mapper_ids_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSocPreprocessRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_soc_preprocess_rules" "test" {
  workspace_id = "%[1]s"
}

locals {
  name       = data.huaweicloud_secmaster_soc_preprocess_rules.test.data[0].name
  mapping_id = data.huaweicloud_secmaster_soc_preprocess_rules.test.data[0].mapping_id
  mapper_id = data.huaweicloud_secmaster_soc_preprocess_rules.test.data[0].mapper_id
}

data "huaweicloud_secmaster_soc_preprocess_rules" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

data "huaweicloud_secmaster_soc_preprocess_rules" "filter_by_mapping_id" {
  workspace_id = "%[1]s"
  mapping_id   = local.mapping_id
}

data "huaweicloud_secmaster_soc_preprocess_rules" "filter_by_mapper_ids" {
  workspace_id = "%[1]s"
  mapper_ids   = [local.mapper_id]
}

output "name_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_preprocess_rules.filter_by_name.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_soc_preprocess_rules.filter_by_name.data[*].name : v == local.name]
  )
}

output "mapping_id_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_preprocess_rules.filter_by_mapping_id.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_soc_preprocess_rules.filter_by_mapping_id.data[*].mapping_id : v == local.mapping_id]
  )
}

output "mapper_ids_filter_useful" {
  value = length(data.huaweicloud_secmaster_soc_preprocess_rules.filter_by_mapper_ids.data) > 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
