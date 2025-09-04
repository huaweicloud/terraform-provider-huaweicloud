package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcGlobalConnectionBandwidthConfigs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_global_connection_bandwidth_configs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcGlobalConnectionBandwidthConfigs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "configs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.bind_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.enable_change_95"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.size_range.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.size_range.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.size_range.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.size_range.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.gcb_type.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.ratio_95peak_plus"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.ratio_95peak_guar"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.sla_level.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.enable_spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.charge_mode.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.crossborder"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.quotas.0.quota"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.quotas.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "configs.0.enable_area_bandwidth"),
				),
			},
		},
	})
}

func testDataSourceCcGlobalConnectionBandwidthConfigs_basic() string {
	return `
data "huaweicloud_cc_global_connection_bandwidth_configs" "test" {}
`
}
