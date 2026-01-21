package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineCheckRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_check_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBaselineCheckRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.severity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.standard"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_type_desc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_rule_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.check_rule_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.scan_result"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.latest_scan_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_num"),
				),
			},
		},
	})
}

func testDataSourceBaselineCheckRules_basic() string {
	return `
data "huaweicloud_hss_baseline_check_rules" "test" {
  type                  = "image"
  image_type            = "registry"
  enterprise_project_id = "all_granted_eps"
  scan_result           = "pass"
}
`
}
