package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscClassificationResults_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_classification_results.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscScanJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscClassificationResults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.asset_id"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.asset_name"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.asset_type"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.object_name"),
					resource.TestCheckResourceAttrSet(dataSource, "classification_list.0.match_info.#"),
					resource.TestCheckOutput("asset_type_filter_is_useful", "true"),
					resource.TestCheckOutput("asset_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDscClassificationResults_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_classification_results" "test" {
  job_id = "%[1]s"
}

# Filter by asset_type
locals {
  asset_type = data.huaweicloud_dsc_classification_results.test.classification_list[0].asset_type
}

data "huaweicloud_dsc_classification_results" "filter_by_asset_type" {
  job_id     = "%[1]s"
  asset_type = local.asset_type
}

locals {
  asset_type_filter_result = [
    for v in data.huaweicloud_dsc_classification_results.filter_by_asset_type.classification_list[*].asset_type : v == local.asset_type
  ]
}

output "asset_type_filter_is_useful" {
  value = alltrue(local.asset_type_filter_result) && length(local.asset_type_filter_result) > 0
}

# Filter by asset_id
locals {
  asset_id = data.huaweicloud_dsc_classification_results.test.classification_list[0].asset_id
}

data "huaweicloud_dsc_classification_results" "filter_by_asset_id" {
  job_id   = "%[1]s"
  asset_id = local.asset_id
}

locals {
  asset_id_filter_result = [
    for v in data.huaweicloud_dsc_classification_results.filter_by_asset_id.classification_list[*].asset_id : v == local.asset_id
  ]
}

output "asset_id_filter_is_useful" {
  value = alltrue(local.asset_id_filter_result) && length(local.asset_id_filter_result) > 0
}
`, acceptance.HW_DSC_SCAN_JOB_ID)
}
